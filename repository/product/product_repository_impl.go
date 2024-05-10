package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
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

func (repository *productRepositoryImpl) Add(ctx context.Context, product product_entity.Product) (*product_entity.Product, error) {
	var productId string
	var createdAt time.Time
	query := `INSERT INTO products
	(
		id, name, sku, category, image_url, notes, price, stock, location, is_available
	)
	VALUES
	(
		gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	RETURNING id, created_at
	`
	if err := repository.dbPool.QueryRow(ctx, query, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable).Scan(&productId, &createdAt); err != nil {
		return &product_entity.Product{}, err
	}

	product.Id = productId
	product.CreatedAt = createdAt.Format(time.RFC3339)

	return &product, nil
}

func (repository *productRepositoryImpl) Delete(ctx context.Context, productId string) (*product_entity.Product, error) {
	query := `update products set is_deleted=true where id=$1 returning id`
	if err := repository.dbPool.QueryRow(ctx, query, productId).Scan(&productId); err != nil {
		return &product_entity.Product{}, err
	}

	product := product_entity.Product{}

	product.Id = productId

	return &product, nil
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

	var dataCount, totalPrice int
	tx, err := repository.dbPool.Begin(ctx)
	if err != nil {
		return &product_entity.ProductCheckout{}, err
	}
	updateProductQuantityQuery := fmt.Sprintf(`
	WITH result as (update products p
	set stock = stock - x.amount 
	from %s
	where p.id = x.id
	returning p.id, p.price, x.amount)
	select count(id) dataCount, SUM(price * amount) totalPrice FROM result
	`, subQuery)
	if err := tx.QueryRow(ctx, updateProductQuantityQuery).Scan(&dataCount, &totalPrice); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	// if one of the product is not exist, the len will be different
	if len(*productCheckout.ProductDetails) != dataCount {
		return &product_entity.ProductCheckout{}, errors.New("no rows in result set")
	}
	// check the paid and change
	if *productCheckout.Paid < totalPrice {
		return &product_entity.ProductCheckout{}, errors.New("paid didn't enough")
	}
	realChange := *productCheckout.Paid - totalPrice
	if realChange != *productCheckout.Change {
		return &product_entity.ProductCheckout{}, errors.New("change is wrong")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

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
	if err := tx.QueryRow(ctx, query, productCheckout.CustomerId, *productCheckout.ProductDetails, productCheckout.Paid, productCheckout.Change).Scan(&productId, &createdAt); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	productCheckout.CheckoutId = productId
	productCheckout.CreatedAt = createdAt.Format(time.RFC3339)

	// updateProductQuantityQuery := fmt.Sprintf(`update products p
	// set stock = stock - x.amount
	// from %s
	// where p.id = x.id
	// `, subQuery)
	// if _, err := tx.Exec(ctx, updateProductQuantityQuery); err != nil {
	// 	return &product_entity.ProductCheckout{}, err
	// }

	return &productCheckout, nil
}

func (repository *productRepositoryImpl) HistorySearch(ctx context.Context, searchQuery product_entity.ProductCheckoutHistoryRequest) ([]product_entity.ProductCheckoutDataResponse, error) {
	query := `SELECT id transactionId, 
		customer_id customerId, 
		product_details productDetails, 
		paid,
		change,
		to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt
		FROM transactions`
	params := []interface{}{}

	if searchQuery.CustomerId != "" {
		query += fmt.Sprintf(" WHERE customer_id = $%s", strconv.Itoa(len(params)+1))
		params = append(params, searchQuery.CustomerId)
	}

	query += " ORDER BY created_at"
	if strings.ToLower(searchQuery.CreatedAt) != "asc" {
		query += " DESC"
	}

	len_p := len(params)
	query += fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(len_p+1), strconv.Itoa(len_p+2))
	params = append(params, searchQuery.Limit, searchQuery.Offset)

	rows, err := repository.dbPool.Query(ctx, query, params...)
	if err != nil {
		return []product_entity.ProductCheckoutDataResponse{}, err
	}
	defer rows.Close()

	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.ProductCheckoutDataResponse])
	if err != nil {
		return []product_entity.ProductCheckoutDataResponse{}, err
	}

	return history, nil

}
