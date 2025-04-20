package lib

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) (map[string]string, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(map[string]string)
		for _, fe := range ve {
			errs[fe.Field()] = HumanMessage(fe)
		}
		return errs, err
	}

	return nil, err
}

func HumanMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
		// case "min":
		// 	return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
		// case "max":
		// 	return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	}
	return fmt.Sprintf("%s is invalid", fe.Field())
}
