package controller

import (
	customer_entity "eniqilo-store/entity/customer"
	exc "eniqilo-store/exceptions"
	auth_service "eniqilo-store/service/auth"
	customer_service "eniqilo-store/service/customer"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	CustomerService customer_service.CustomerService
	AuthService     auth_service.AuthService
}

func NewCustomerController(customerService customer_service.CustomerService) CustomerController {
	return CustomerController{
		CustomerService: customerService,
	}
}

func (controller *CustomerController) Register(ctx *fiber.Ctx) error {
	customerReq := new(customer_entity.CustomerRegisterRequest)
	if err := ctx.BodyParser(customerReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.CustomerService.Register(ctx.UserContext(), *customerReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller *CustomerController) Search(ctx *fiber.Ctx) error {
	customerQuery := new(customer_entity.CustomerSearchRequest)
	if err := ctx.QueryParser(customerQuery); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.CustomerService.Search(ctx.UserContext(), *customerQuery)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
