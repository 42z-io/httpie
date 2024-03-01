package httpie

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteErrJson(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteErrJson(w, http.StatusNotFound, "not found")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"not found"}`, w.Body.String())
}

func TestWriteErrHttpErr(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteErr(w, ErrBadRequest)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"bad request"}`, w.Body.String())
}

func TestWriteErrValidationErr(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteErr(w, NewErrHttpValidation(map[string]string{"name": "required"}))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"errors":{"name":"required"},"message":"validation failed"}`, w.Body.String())
}

func TestWriteErrOther(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteErr(w, errors.New("test"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"internal server error"}`, w.Body.String())
}

func TestWriteOk(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteOk(w, map[string]string{"message": "hello"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"hello"}`, w.Body.String())
}

func TestWriteAccepted(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	WriteAccepted(w)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Header().Get("Content-Type"))
	assert.Empty(t, w.Body.String())
}

func TestWriteOkOrErrErr(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	var data map[string]string
	WriteOkOrErr(w, data, errors.New("test"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"internal server error"}`, w.Body.String())
}

func TestWriteOkOrErrOk(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	data := map[string]string{"name": "hello"}
	WriteOkOrErr(w, data, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"name":"hello"}`, w.Body.String())
}
