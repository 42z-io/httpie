package httpie

import (
	"net/http"
)

type IErrHttp interface {
	StatusCode() int
	Error() string
}

type IErrHttpValidation interface {
	StatusCode() int
	Error() string
	ValidationErrors() map[string]string
}

type ErrHttp struct {
	statusCode int
	error      string
}

func (e ErrHttp) StatusCode() int {
	return e.statusCode
}

func (e ErrHttp) Error() string {
	return e.error
}

func NewErrHttp(statusCode int, error string) ErrHttp {
	return ErrHttp{statusCode, error}
}

type ErrHttpValidation struct {
	validationErrors map[string]string
	ErrHttp
}

func (e ErrHttpValidation) ValidationErrors() map[string]string {
	return e.validationErrors
}

func NewErrHttpValidation(errors map[string]string) ErrHttpValidation {
	return ErrHttpValidation{errors, ErrHttp{http.StatusBadRequest, "validation failed"}}
}

// These are standard errors that should be returned at the repository level, its not meant to be
// exhaustive of all HTTP errors but rather standard ones that make sense to propagate up from the services and repositories.
var (
	ErrNotFound     = NewErrHttp(http.StatusNotFound, "not found")
	ErrUnauthorized = NewErrHttp(http.StatusUnauthorized, "unauthorized")
	ErrBadRequest   = NewErrHttp(http.StatusBadRequest, "bad request")
	ErrForbidden    = NewErrHttp(http.StatusForbidden, "forbidden")
	ErrConflict     = NewErrHttp(http.StatusConflict, "conflict")
	ErrInternal     = NewErrHttp(http.StatusInternalServerError, "internal server error")
)
