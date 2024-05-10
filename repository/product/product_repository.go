package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
)

type ProductRepository interface {
	Add(ctx context.Context, req product_entity.Product) (*product_entity.Product, error)
	Edit(ctx context.Context, req product_entity.Product, productId string) (*product_entity.Product, error)
	Search(ctx context.Context, req product_entity.ProductSearch) (*[]product_entity.ProductSearchData, error)
	Delete(ctx context.Context, catId string) (*product_entity.ProductDeleteData, error)
	Checkout(ctx context.Context, req product_entity.ProductCheckout) (*product_entity.ProductCheckout, error)
	HistorySearch(ctx context.Context, req product_entity.ProductCheckoutHistoryRequest) ([]product_entity.ProductCheckoutDataResponse, error)
	CustomerSearch(ctx context.Context, query product_entity.ProductCustomerSearch) (*[]product_entity.ProductCustomerSearchData, error)
}
