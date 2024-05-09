package product_service

import (
	product_entity "eniqilo-store/entity/product"

	"github.com/gofiber/fiber/v2"
)

type ProductService interface {
	Add(ctx *fiber.Ctx, req product_entity.ProductRegisterRequest) (product_entity.ProductRegisterResponse, error)
	Edit(ctx *fiber.Ctx, req product_entity.ProductEditRequest) (product_entity.ProductEditResponse, error)
	Search(ctx *fiber.Ctx, searchQueries product_entity.ProductSearchQuery) (product_entity.ProductSearchResponse, error)
}
