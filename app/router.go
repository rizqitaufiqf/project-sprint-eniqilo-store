package app

import (
	"eniqilo-store/controller"
	"eniqilo-store/helpers"

	product_repository "eniqilo-store/repository/product"
	staffr_repository "eniqilo-store/repository/staff"
	auth_service "eniqilo-store/service/auth"
	product_service "eniqilo-store/service/product"
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

	productRepository := product_repository.NewProductRepository()
	productService := product_service.NewProductService(productRepository, dbPool, authService, validator)
	productController := controller.NewProductController(productService, authService)

	// Staffs API
	staffApi := app.Group("/v1/staff")
	staffApi.Post("/register", staffController.Register)
	staffApi.Post("/login", staffController.Login)

	// JWT middleware
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	// Products API
	productApi := app.Group("/v1/product")
	productApi.Post("/", productController.Add)
}
