package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepositoryImpl struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) ProductRepository {
	return &productRepositoryImpl{
		dbPool: dbPool,
	}
}

func (repository *productRepositoryImpl) Add(ctx context.Context, product product_entity.ProductRegisterRequest) (*product_entity.ProductData, error) {
	var productId, createdAt string
	query := `INSERT INTO products
	(
		name, sku, category, image_url, notes, price, stock, location, is_available
	)
	VALUES
	(
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	RETURNING id, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt
	`
	if err := repository.dbPool.QueryRow(ctx, query, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable).Scan(&productId, &createdAt); err != nil {
		return &product_entity.ProductData{}, err
	}

	return &product_entity.ProductData{
		Id:        productId,
		CreatedAt: createdAt,
	}, nil
}

func (repository *productRepositoryImpl) Edit(ctx context.Context, product product_entity.ProductEditRequest, productId string) (*product_entity.ProductData, error) {
	updateQ := `UPDATE products SET name = $1, sku = $2, category = $3,
	image_url = $4, notes = $5, price = $6, stock = $7, location = $8,
	is_available = $9
	WHERE id = $10
	`
	res, err := repository.dbPool.Exec(ctx, updateQ, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, productId)
	if err != nil {
		return &product_entity.ProductData{}, err
	}
	if res.RowsAffected() == 0 {
		return &product_entity.ProductData{}, exc.NotFoundException("Product id does not exist")
	}

	return &product_entity.ProductData{
		Id: productId,
	}, nil
}

func (repository *productRepositoryImpl) Delete(ctx context.Context, productId string) (*product_entity.ProductDeleteData, error) {
	query := `update products set is_deleted=true where id=$1 AND is_deleted = false returning id`
	if err := repository.dbPool.QueryRow(ctx, query, productId).Scan(&productId); err != nil {
		return &product_entity.ProductDeleteData{}, err
	}

	product := product_entity.ProductDeleteData{}

	product.Id = productId

	return &product, nil
}

func (repository *productRepositoryImpl) Search(ctx context.Context, searchQuery product_entity.ProductSearchQuery) (*[]product_entity.ProductSearchData, error) {
	query := `SELECT id, name, sku, category, image_url, stock, notes, price, location, is_available, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt FROM products WHERE is_deleted = false`
	var whereClause []string
	var searchParams []interface{}

	if searchQuery.Id != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Id)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.IsAvailable != nil {
		whereClause = append(whereClause, fmt.Sprintf("is_available = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.IsAvailable)
	}
	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Category)
	}
	if searchQuery.Sku != "" {
		whereClause = append(whereClause, fmt.Sprintf("sku = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Sku)
	}
	if searchQuery.InStock {
		query += "stock > 0"
	}

	if len(whereClause) > 0 {
		query += " AND " + strings.Join(whereClause, " AND ")
	}

	// construct order by
	var orderByClause []string
	var orderByDefault = ` ORDER BY created_at DESC`
	if strings.ToLower(searchQuery.Price) == "asc" {
		query += "price asc"
	} else if strings.ToLower(searchQuery.Price) == "desc" {
		query += "price desc"
	}

	if strings.ToLower(searchQuery.CreatedAt) == "asc" {
		orderByClause = append(orderByClause, "created_at")
	} else if strings.ToLower(searchQuery.CreatedAt) == "desc" {
		orderByClause = append(orderByClause, "created_at desc")
	}

	if len(orderByClause) > 0 {
		query += " ORDER BY " + strings.Join(orderByClause, ", ")
	} else {
		query += orderByDefault
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", searchQuery.Limit, searchQuery.Offset)
	fmt.Println(query, searchParams)
	rows, err := repository.dbPool.Query(ctx, query, searchParams...)
	if err != nil {
		fmt.Println(err.Error())
		return &[]product_entity.ProductSearchData{}, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.ProductSearchData])
	if err != nil {
		return &[]product_entity.ProductSearchData{}, err
	}

	return &products, nil
}

func (repository *productRepositoryImpl) Checkout(ctx context.Context, productCheckout product_entity.ProductCheckout) (*product_entity.ProductCheckout, error) {
	// first we make sure that the customer id is there
	var cId string
	cIdQuery := `select id from customers where id = $1`
	if err := repository.dbPool.QueryRow(ctx, cIdQuery, productCheckout.CustomerId).Scan(&cId); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	var updatedQuantityId []string
	for _, item := range *productCheckout.ProductDetails {
		x := fmt.Sprintf(`('%s',%d)`, item.ProductId, item.Quantity)
		updatedQuantityId = append(updatedQuantityId, x)
	}
	subQuery := fmt.Sprintf(`(values %s) as x(id, amount)`, strings.Join(updatedQuantityId, ", "))

	var dataCount, totalPrice, unavailableCount, outOfStock *int
	productPriceQuery := fmt.Sprintf(`
	SELECT count(p.id) dataCount, 
		SUM(p.price * x.amount) totalPrice,
		COUNT(CASE WHEN is_available = false THEN 1 END) AS unavailableCount,
		COUNT(CASE WHEN stock < x.amount THEN 1 END) AS outOfStock
	FROM products p join %s 
		ON x.id = p.id `, subQuery)
	if err := repository.dbPool.QueryRow(ctx, productPriceQuery).Scan(&dataCount, &totalPrice, &unavailableCount, &outOfStock); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	// if one of the product is not exist, the len will be different
	if len(*productCheckout.ProductDetails) != *dataCount {
		err := errors.New("no rows in result set")
		return &product_entity.ProductCheckout{}, err
	}

	// one of products is not available
	if *unavailableCount > 0 {
		err := errors.New("doesn’t pass validation: one of productIds is not available")
		return &product_entity.ProductCheckout{}, err
	}

	// outOfStock
	if *outOfStock > 0 {
		err := errors.New("doesn’t pass validation: out of stock")
		return &product_entity.ProductCheckout{}, err
	}

	// check the paid and change
	if *productCheckout.Paid < *totalPrice {
		err := errors.New("doesn’t pass validation: paid didn't enough")
		return &product_entity.ProductCheckout{}, err
	}
	realChange := *productCheckout.Paid - *totalPrice
	if realChange != *productCheckout.Change {
		err := errors.New("doesn’t pass validation: change is wrong")
		return &product_entity.ProductCheckout{}, err
	}

	updateProductQuantityQuery := fmt.Sprintf(`update products p
	set stock = stock - x.amount 
	from %s
	where p.id = x.id
	`, subQuery)
	if _, err := repository.dbPool.Exec(ctx, updateProductQuantityQuery); err != nil {
		return &product_entity.ProductCheckout{}, err
	}
	var productId string
	var createdAt time.Time
	query := `INSERT INTO transactions
	(
		id, customer_id, product_details, paid, change
	)
	VALUES
	(
		gen_random_uuid(), $1, $2, $3, $4
	)
	RETURNING id, created_at
	`
	if err := repository.dbPool.QueryRow(ctx, query, productCheckout.CustomerId, *productCheckout.ProductDetails, productCheckout.Paid, productCheckout.Change).Scan(&productId, &createdAt); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	productCheckout.CheckoutId = productId
	productCheckout.CreatedAt = createdAt.Format(time.RFC3339)

	return &productCheckout, nil
}

func (repository *productRepositoryImpl) HistorySearch(ctx context.Context, searchQuery product_entity.ProductCheckoutHistoryRequest) (*[]product_entity.ProductCheckoutDataResponse, error) {
	query := `SELECT id transactionId, 
		customer_id customerId, 
		product_details productDetails, 
		paid,
		change,
		to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt
		FROM transactions`
	params := []interface{}{}

	if searchQuery.CustomerId != "" {
		query += fmt.Sprintf(" WHERE customer_id = $%d", len(params)+1)
		params = append(params, searchQuery.CustomerId)
	}

	query += " ORDER BY created_at"
	if strings.ToLower(searchQuery.CreatedAt) != "asc" {
		query += " DESC"
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", searchQuery.Limit, searchQuery.Offset)

	rows, err := repository.dbPool.Query(ctx, query, params...)
	if err != nil {
		return &[]product_entity.ProductCheckoutDataResponse{}, err
	}
	defer rows.Close()

	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.ProductCheckoutDataResponse])
	if err != nil {
		return &[]product_entity.ProductCheckoutDataResponse{}, err
	}

	return &history, nil
}

func (repository *productRepositoryImpl) CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearchQuery) (*[]product_entity.ProductCustomerSearchData, error) {
	query := `SELECT id, name, sku, category, image_url, stock, price, location, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt FROM products WHERE is_deleted = FALSE AND is_available = TRUE`
	var whereClause []string
	params := []interface{}{}

	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Name)
	}

	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Category)
	}

	if searchQuery.Sku != "" {
		whereClause = append(whereClause, fmt.Sprintf("sku = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Sku)
	}

	if searchQuery.InStock {
		whereClause = append(whereClause, "stock > 0")
	}

	if len(whereClause) > 0 {
		query += " AND " + strings.Join(whereClause, " AND ")
	}

	if searchQuery.Price == "asc" || searchQuery.Price == "desc" {
		query += fmt.Sprintf(" ORDER BY price %s", searchQuery.Price)
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", searchQuery.Limit, searchQuery.Offset)

	rows, err := repository.dbPool.Query(ctx, query, params...)
	if err != nil {
		return &[]product_entity.ProductCustomerSearchData{}, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.ProductCustomerSearchData])
	if err != nil {
		return &[]product_entity.ProductCustomerSearchData{}, err
	}

	return &products, nil
}
