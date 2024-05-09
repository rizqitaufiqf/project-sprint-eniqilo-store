package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	"time"

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

func (repository *productRepositoryImpl) Edit(ctx context.Context, product product_entity.Product, productId string) (*product_entity.Product, error) {
	updateQ := `UPDATE products SET name = $1, sku = $2, category = $3,
	image_url = $4, notes = $5, price = $6, stock = $7, location = $8,
	is_available = $9
	WHERE id = $10
	`
	res, err := repository.dbPool.Exec(ctx, updateQ, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, productId)
	if err != nil {
		return &product_entity.Product{}, err
	}
	if res.RowsAffected() == 0 {
		return &product_entity.Product{}, exc.NotFoundException("Product id does not exist")
	}

	product.Id = productId

	return &product, nil
}
