package httpie

import (
	"net/http"
)

type HttpError interface {
	StatusCode() int
	Error() string
}

type HttpErrorImpl struct {
	statusCode int
	error      string
}

func (e HttpErrorImpl) StatusCode() int {
	return e.statusCode
}

func (e HttpErrorImpl) Error() string {
	return e.error
}

func NewHttpError(statusCode int, error string) HttpErrorImpl {
	return HttpErrorImpl{statusCode, error}
}

// These are standard errors that should be returned at the repository level, its not meant to be
// exhaustive of all HTTP errors but rather standard ones that make sense to propagate up from the services and repositories.
var (
	ErrNotFound     = NewHttpError(http.StatusNotFound, "not found")
	ErrUnauthorized = NewHttpError(http.StatusUnauthorized, "unauthorized")
	ErrBadRequest   = NewHttpError(http.StatusBadRequest, "bad request")
	ErrForbidden    = NewHttpError(http.StatusForbidden, "forbidden")
	ErrConflict     = NewHttpError(http.StatusConflict, "conflict")
	ErrInternal     = NewHttpError(http.StatusInternalServerError, "internal server error")
)
