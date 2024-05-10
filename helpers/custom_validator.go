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

func validateUrl(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	pattern := `^(?:https?:\/\/)?(?:www\.)?(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?:\/[^\s]*)?$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func validateProductCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	for _, categ := range ProductCategory {
		if categ == value {
			return true
		}
	}
	return false
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("phoneNumber", validatePhoneNumber)
	validator.RegisterValidation("productCategory", validateProductCategory)
	validator.RegisterValidation("validateUrl", validateUrl)
}
