package httpie

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
}

func TestReadJsonOk(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "http://example.com", strings.NewReader(`{"name":"hello"}`))
	var data testStruct
	err := ReadJson(r, &data)
	assert.Nil(t, err)
	assert.Equal(t, "hello", data.Name)
}

func TestReadJsonErrSyntax(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "http://example.com", strings.NewReader(`{"name`))
	var data testStruct
	err := ReadJson(r, &data)
	assert.Error(t, err)
}

func TestReadJsonErrIO(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("POST", "http://example.com", iotest.ErrReader(errors.New("io error")))
	var data testStruct
	err := ReadJson(r, &data)
	assert.Error(t, err)
}
