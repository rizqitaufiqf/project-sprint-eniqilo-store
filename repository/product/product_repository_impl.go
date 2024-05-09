package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	"errors"
	"fmt"
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

	var productIds []string
	idToQuantity := make(map[string]int)
	for _, item := range *productCheckout.ProductDetails {
		id := item.ProductId
		productIds = append(productIds, id)
		idToQuantity[item.ProductId] = item.Quantity
	}
	productPriceQuery := `select id, price, stock, is_available from products where id = any($1)`
	rows, err := repository.dbPool.Query(ctx, productPriceQuery, productIds)
	if err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	type ProductPrice struct {
		Id          string
		Price       int
		Stock       int
		IsAvailable bool
	}
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProductPrice])
	if err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	// if one of the product is not exist, the len will be different
	if len(products) != len(productIds) {
		return &product_entity.ProductCheckout{}, errors.New("no rows in result set")
	}
	var totalPrice int
	var updatedQuantityId string

	for _, item := range products {
		isAvailable := item.IsAvailable
		if !isAvailable {
			return &product_entity.ProductCheckout{}, errors.New("one of the product isn't available")
		}
		if item.Stock < idToQuantity[item.Id] {
			return &product_entity.ProductCheckout{}, errors.New("stock didn't enough")
		}

		totalPrice += item.Price
		updatedQuantityId += fmt.Sprintf("when id = '%s' then %d \n", item.Id, item.Stock-idToQuantity[item.Id])
	}

	// check the paid and change
	if productCheckout.Paid < totalPrice {
		return &product_entity.ProductCheckout{}, errors.New("paid didn't enough")
	}
	realChange := productCheckout.Paid - totalPrice
	if realChange != productCheckout.Change {
		return &product_entity.ProductCheckout{}, errors.New("change is wrong")
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
	fmt.Println(fmt.Sprintf(`update products
	set stock = 
		case
		%s
		end
	returning id
	`, updatedQuantityId))

	var tmpProductId string
	updateProductQuantityQuery := fmt.Sprintf(`update products
	set stock = 
		case
		%s
		end
	returning id
	`, updatedQuantityId)
	if err := repository.dbPool.QueryRow(ctx, updateProductQuantityQuery).Scan(&tmpProductId); err != nil {
		return &product_entity.ProductCheckout{}, err
	}

	return &productCheckout, nil
}
