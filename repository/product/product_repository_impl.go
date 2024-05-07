package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) Add(ctx context.Context, dbPool *pgxpool.Pool, product product_entity.Product) (*product_entity.Product, error) {
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
	if err := dbPool.QueryRow(ctx, query, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable).Scan(&productId, &createdAt); err != nil {
		return &product_entity.Product{}, err
	}

	product.Id = productId
	product.CreatedAt = createdAt.Format(time.RFC3339)

	return &product, nil
}
