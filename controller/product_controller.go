package controller

import (
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	auth_service "eniqilo-store/service/auth"
	product_service "eniqilo-store/service/product"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService product_service.ProductService
	AuthService    auth_service.AuthService
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
