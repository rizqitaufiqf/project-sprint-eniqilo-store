package helpers

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator"
)

func validateBoolean(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := strconv.ParseBool(value)

	return err == nil
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	pattern := `^\+\d{1,}(?:-?\d{1,})+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("boolean", validateBoolean)
	validator.RegisterValidation("phoneNumber", validatePhoneNumber)
}
