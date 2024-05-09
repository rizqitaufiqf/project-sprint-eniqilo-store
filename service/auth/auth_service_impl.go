package auth_service

import (
	"context"
	helpers "eniqilo-store/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// 8 hours
var expDuration = time.Now().Add(time.Hour * 8).Unix()

type authServiceImpl struct {
}

func NewAuthService() AuthService {
	return &authServiceImpl{}
}

func (service *authServiceImpl) GenerateToken(ctx context.Context, userId string) (string, error) {
	jwtconf := jwt.MapClaims{
		"user_id": userId,
		"exp":     expDuration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtconf)
	signToken, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signToken, nil
}

func (service *authServiceImpl) GetValidUser(ctx *fiber.Ctx) (string, error) {
	userInfo := ctx.Locals(helpers.JwtContextKey).(*jwt.Token)
	// convert userInfo claims to jwt mapclaims
	jwtconf := userInfo.Claims.((jwt.MapClaims))
	// convert user_id to string
	userId := jwtconf["user_id"].(string)

	return userId, nil
}
