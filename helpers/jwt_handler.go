package helpers

import (
	exc "eniqilo-store/exceptions"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func GetTokenHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(viper.GetString("JWT_SECRET"))},
		ContextKey: JwtContextKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
	})
}

func CheckTokenHeader(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return exc.UnauthorizedException("Unauthorized")
	} else {
		return ctx.Next()
	}
}
