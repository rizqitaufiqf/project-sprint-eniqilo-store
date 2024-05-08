package app

import (
	"eniqilo-store/controller"
	"eniqilo-store/helpers"

	customer_repository "eniqilo-store/repository/customer"
	staffr_repository "eniqilo-store/repository/staff"
	auth_service "eniqilo-store/service/auth"
	customer_service "eniqilo-store/service/customer"
	staff_service "eniqilo-store/service/staff"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBluePrint(app *fiber.App, dbPool *pgxpool.Pool) {
	validator := validator.New()
	// register custom validator
	helpers.RegisterCustomValidator(validator)

	authService := auth_service.NewAuthService()

	staffRepository := staffr_repository.NewStaffRepository()
	staffService := staff_service.NewStaffService(staffRepository, dbPool, authService, validator)
	staffController := controller.NewStaffController(staffService, authService)

	customerRepository := customer_repository.NewCustomerRepository(dbPool)
	customerService := customer_service.NewCustomerService(customerRepository, validator)
	customerController := controller.NewCustomerController(customerService)

	// Staffs API
	staffApi := app.Group("/v1/staff")
	staffApi.Post("/register", staffController.Register)
	staffApi.Post("/login", staffController.Login)

	// JWT middleware
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	// Staffs API
	customerApi := app.Group("/v1/customer")
	customerApi.Post("/register", customerController.Register)
	customerApi.Get("/", customerController.Search)
}
