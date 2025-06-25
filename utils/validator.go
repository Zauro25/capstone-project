package utils

import (
	"github.com/go-playground/validator/v10"
)

// Validator instance
var Validate = validator.New()

// ValidateStruct memvalidasi struct berdasarkan tag 'validate'
func ValidateStruct(s interface{}) error {
	return Validate.Struct(s)
}
