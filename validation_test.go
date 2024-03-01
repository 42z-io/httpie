package httpie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValidateStruct struct {
	Name     string `json:"other" validate:"required"`
	Password string `validate:"required,securepassword"`
	Age      int    `validate:"required,max=100"`
}

func TestValidatorOk(t *testing.T) {
	err := Validate(testValidateStruct{Name: "test", Password: "test.test.test12345678", Age: 50})
	assert.Nil(t, err)
}

func TestValidatorErr(t *testing.T) {
	err := Validate(testValidateStruct{Name: "test"})
	assert.NotNil(t, err)
	assert.Equal(t, "required", err.ValidationErrors()["Password"])
}

func TestValidatorErrWithParam(t *testing.T) {
	err := Validate(testValidateStruct{Name: "test", Age: 2000})
	assert.NotNil(t, err)
	assert.Equal(t, "max=100", err.ValidationErrors()["Age"])
}

func TestValidatorJsonFieldErr(t *testing.T) {
	err := Validate(testValidateStruct{})
	assert.NotNil(t, err)
	assert.Equal(t, "required", err.ValidationErrors()["other"])
}

func TestValidatorErrSecurePassword(t *testing.T) {
	err := Validate(testValidateStruct{Name: "test", Password: "test"})
	assert.NotNil(t, err)
	assert.Equal(t, "securepassword", err.ValidationErrors()["Password"])
}
