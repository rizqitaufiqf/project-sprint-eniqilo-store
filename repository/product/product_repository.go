package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
)

type ProductRepository interface {
	Add(ctx context.Context, req product_entity.Product) (*product_entity.Product, error)
	CustomerSearch(ctx context.Context, query product_entity.ProductCustomerSearch) (*[]product_entity.ProductCustomerSearchData, error)
}
