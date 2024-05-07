package app

import (
	"eniqilo-store/helpers"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBluePrint(app *fiber.App, dbPool *pgxpool.Pool) {
	validator := validator.New()
	// register custom validator
	helpers.RegisterCustomValidator(validator)

	// Users API
	userApi := app.Group("/v1/")
	userApi.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("not protected")
	})

	// JWT middleware
	app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())

	userApi.Get("/hello/protected", func(c *fiber.Ctx) error {
		return c.SendString("protected")
	})
}
