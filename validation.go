package httpie

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
	Validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return field.Name
		}
		return name
	})
	Validator.RegisterValidation("securepassword", SecurePasswordValidator)
}

// Validate an object and return a validation error if any
func Validate(s any) IErrHttpValidation {
	var validations = map[string]string{}
	err := Validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.Tag()
			param := err.Param()
			var message string = tag
			if param != "" {
				message += "=" + param
			}
			validations[err.Field()] = message
		}
		return NewErrHttpValidation(validations)
	}
	return nil
}

const MIN_ENTROPY = 60

// Custom validator for secure passwords - ensures entropy is greater than 60
func SecurePasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	err := passwordvalidator.Validate(password, MIN_ENTROPY)
	return err == nil
}
