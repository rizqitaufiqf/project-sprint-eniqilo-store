package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
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

func (repository *productRepositoryImpl) Delete(ctx context.Context, productId string) (*product_entity.Product, error) {
	query := `update products set is_deleted=true where id=$1 returning id`
	if err := repository.dbPool.QueryRow(ctx, query, productId).Scan(&productId); err != nil {
		return &product_entity.Product{}, err
	}

	product := product_entity.Product{}

	product.Id = productId

	return &product, nil
}
