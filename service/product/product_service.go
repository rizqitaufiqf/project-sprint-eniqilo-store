package product_service

import (
	"context"
	product_entity "eniqilo-store/entity/product"

	"github.com/gofiber/fiber/v2"
)

type ProductService interface {
	Add(ctx *fiber.Ctx, req product_entity.ProductRegisterRequest) (product_entity.ProductRegisterResponse, error)
	CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearchQuery) (product_entity.ProductCustomerSearchResponse, error)
}
