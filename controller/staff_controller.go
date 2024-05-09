package controller

import (
	staff_entity "eniqilo-store/entity/staff"
	exc "eniqilo-store/exceptions"
	staff_service "eniqilo-store/service/staff"

	"github.com/gofiber/fiber/v2"
)

type StaffController struct {
	StaffService staff_service.StaffService
}

func NewStaffController(staffService staff_service.StaffService) *StaffController {
	return &StaffController{
		StaffService: staffService,
	}
}

func (controller StaffController) Register(ctx *fiber.Ctx) error {
	staffReq := new(staff_entity.StaffRegisterRequest)
	if err := ctx.BodyParser(staffReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.StaffService.Register(ctx.UserContext(), *staffReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller StaffController) Login(ctx *fiber.Ctx) error {
	staffReq := new(staff_entity.StaffLoginRequest)
	if err := ctx.BodyParser(staffReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.StaffService.Login(ctx.UserContext(), *staffReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
