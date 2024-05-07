package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Add(ctx context.Context, dbPool *pgxpool.Pool, req product_entity.Product) (*product_entity.Product, error)
}
