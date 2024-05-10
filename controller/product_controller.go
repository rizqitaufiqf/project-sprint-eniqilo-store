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
func (controller *ProductController) Delete(ctx *fiber.Ctx) error {
	resp, err := controller.ProductService.Delete(ctx)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)

}

func (controller *ProductController) Checkout(ctx *fiber.Ctx) error {
	productReq := new(product_entity.ProductCheckoutRequest)
	if err := ctx.BodyParser(productReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.ProductService.Checkout(ctx, *productReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller *ProductController) History(ctx *fiber.Ctx) error {
	historySearchQuery := new(product_entity.ProductCheckoutHistoryRequest)
	historySearchQuery.Limit = 5
	historySearchQuery.Offset = 0
	if err := ctx.QueryParser(historySearchQuery); err != nil {
		return exc.BadRequestException("Error when parsing request query")
	}

	resp, err := controller.ProductService.HistorySearch(ctx, *historySearchQuery)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)

}
