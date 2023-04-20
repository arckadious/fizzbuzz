package validator

import (
	"github.com/go-playground/validator/v10"
)

// New create Config
func New() *validator.Validate {

	validate := validator.New()

	// validate.RegisterValidation("multiples", func(fl validator.FieldLevel) bool {

	// 	return true
	// })

	return validate
}
