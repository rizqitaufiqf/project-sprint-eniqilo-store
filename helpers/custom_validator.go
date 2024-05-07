package helpers

import (
	"regexp"

	"github.com/go-playground/validator"
)

func validatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	pattern := `^\+\d{1,}(?:-?\d{1,})+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("phoneNumber", validatePhoneNumber)
}
