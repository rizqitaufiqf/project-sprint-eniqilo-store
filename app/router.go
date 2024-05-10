package app

import (
	"eniqilo-store/controller"
	"eniqilo-store/helpers"

	product_repository "eniqilo-store/repository/product"
	staff_repository "eniqilo-store/repository/staff"
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

	staffRepository := staff_repository.NewStaffRepository(dbPool)
	staffService := staff_service.NewStaffService(staffRepository, authService, validator)
	staffController := controller.NewStaffController(staffService)

	productRepository := product_repository.NewProductRepository(dbPool)
	productService := product_service.NewProductService(productRepository, validator)
	productController := controller.NewProductController(productService)

	// Staffs API
	staffApi := app.Group("/v1/staff")
	staffApi.Post("/register", staffController.Register)
	staffApi.Post("/login", staffController.Login)

	//Search SKU
	app.Get("/v1/product/customer", productController.CustomerSearch)

	// JWT middleware
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	// Products API
	productApi := app.Group("/v1/product")
	productApi.Post("/", productController.Add)
}
