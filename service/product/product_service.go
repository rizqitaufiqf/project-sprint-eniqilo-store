package product_service

import (
	"context"
	product_entity "eniqilo-store/entity/product"

	"github.com/gofiber/fiber/v2"
)

type ProductService interface {
	Add(ctx *fiber.Ctx, req product_entity.ProductRegisterRequest) (product_entity.ProductRegisterResponse, error)
	Delete(ctx *fiber.Ctx) (product_entity.ProductDeleteResponse, error)
	Checkout(ctx *fiber.Ctx, req product_entity.ProductCheckoutRequest) (product_entity.ProductCheckoutResponse, error)
	HistorySearch(ctx *fiber.Ctx, req product_entity.ProductCheckoutHistoryRequest) (product_entity.ProductCheckoutHistoryResponse, error)
	CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearchQuery) (product_entity.ProductCustomerSearchResponse, error)
}
