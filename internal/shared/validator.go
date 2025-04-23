// lib/validator.go
package shared

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
	Validate.RegisterValidation("not_pending", func(fl validator.FieldLevel) bool {
		return fl.Field().String() != "pending"
	})
}
