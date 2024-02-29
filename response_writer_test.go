package httpie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatchedResponseWriter(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	w := NewWatchedResponseWriter(rr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
	w.Apply()
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.StatusOK, w.StatusCode())
	assert.Equal(t, 5, w.BytesWritten())
	assert.Equal(t, 5, len(rr.Body.String()))
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestWatchedResponseWriterReset(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	w := NewWatchedResponseWriter(rr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
	w.Reset()
	w.WriteHeader(http.StatusNotFound)
	w.Apply()
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, 0, w.BytesWritten())
	assert.Equal(t, 0, len(rr.Body.String()))
}
