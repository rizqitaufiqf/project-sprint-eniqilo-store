package app

import (
	"eniqilo-store/controller"
	"eniqilo-store/helpers"

	customer_repository "eniqilo-store/repository/customer"
	product_repository "eniqilo-store/repository/product"
	staff_repository "eniqilo-store/repository/staff"
	auth_service "eniqilo-store/service/auth"
	customer_service "eniqilo-store/service/customer"
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

	// Customer API
	customerApi := app.Group("/v1/customer")
	customerApi.Post("/register", customerController.Register)
	customerApi.Get("/", customerController.Search)

	// Products API
	productApi := app.Group("/v1/product")
	productApi.Post("/", productController.Add)
	productApi.Delete("/:id", productController.Delete)
}
