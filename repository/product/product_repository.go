package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
)

type ProductRepository interface {
	Add(ctx context.Context, req product_entity.Product) (*product_entity.Product, error)
	Delete(ctx context.Context, productId string) (*product_entity.Product, error)
}
