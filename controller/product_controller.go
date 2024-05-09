package controller

import (
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	product_service "eniqilo-store/service/product"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService product_service.ProductService
}

func NewProductController(productService product_service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

func (controller *ProductController) Add(ctx *fiber.Ctx) error {
	productReq := new(product_entity.ProductRegisterRequest)
	if err := ctx.BodyParser(productReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.ProductService.Add(ctx, *productReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller *ProductController) Edit(ctx *fiber.Ctx) error {
	editProductReq := new(product_entity.ProductEditRequest)
	if err := ctx.BodyParser(editProductReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.ProductService.Edit(ctx, *editProductReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller *ProductController) Search(ctx *fiber.Ctx) error {
	productSearchQueries := new(product_entity.ProductSearchQuery)
	if err := ctx.QueryParser(productSearchQueries); err != nil {
		return exc.BadRequestException("Error when parsing request query")
	}

	resp, err := controller.ProductService.Search(ctx, *productSearchQueries)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
