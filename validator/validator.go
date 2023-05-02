// This package creates a custom validator
package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator class
type Validator *validator.Validate

// New construtor Validator
func New() Validator {

	validate := validator.New()

	// validate.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
	// 	return true
	// })

	return validate
}
