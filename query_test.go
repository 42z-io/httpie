package httpie

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryParamDefaultExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?name=hello", nil)
	value := GetQueryParamDefault(r, "name", "world")
	assert.Equal(t, "hello", value)
}

func TestGetQueryParamDefaultNotExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	value := GetQueryParamDefault(r, "name", "world")
	assert.Equal(t, "world", value)
}

func TestGetQueryParamExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?name=hello", nil)
	ok, value := GetQueryParam(r, "name")
	assert.True(t, ok)
	assert.Equal(t, "hello", value)
}

func TestGetQueryParamNotExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	ok, value := GetQueryParam(r, "name")
	assert.False(t, ok)
	assert.Equal(t, "", value)
}

func TestGetQueryParamIntDefaultExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?age=10", nil)
	value, err := GetQueryParamIntDefault(r, "age", 20)
	assert.NoError(t, err)
	assert.Equal(t, 10, value)
}

func TestGetQueryParamIntDefaultNotExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	value, err := GetQueryParamIntDefault(r, "age", 20)
	assert.NoError(t, err)
	assert.Equal(t, 20, value)
}

func TestGetQueryParamIntDefaultError(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?age=hello", nil)
	value, err := GetQueryParamIntDefault(r, "age", 22)
	assert.Equal(t, 0, value)
	assert.Error(t, err)
}

func TestGetQueryParamIntError(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?age=hello", nil)
	ok, value, err := GetQueryParamInt(r, "age")
	assert.True(t, ok)
	assert.Equal(t, 0, value)
	assert.Error(t, err)
}

func TestGetQueryParamListDefaultExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?name=hello,hi", nil)
	value := GetQueryParamListDefault(r, "name", []string{"world"})
	assert.Equal(t, []string{"hello", "hi"}, value)
}

func TestGetQueryParamListDefaultNotExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	value := GetQueryParamListDefault(r, "name", []string{"world"})
	assert.Equal(t, []string{"world"}, value)
}

func TestGetQueryParamListNotExists(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com", nil)
	ok, value := GetQueryParamList(r, "name")
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestGetQueryParamListExistsMany(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?name=1,2,3", nil)
	ok, value := GetQueryParamList(r, "name")
	assert.True(t, ok)
	assert.Equal(t, []string{"1", "2", "3"}, value)
}

func TestGetQueryParamListExistsOne(t *testing.T) {
	t.Parallel()
	r := httptest.NewRequest("GET", "http://example.com?name=1", nil)
	ok, value := GetQueryParamList(r, "name")
	assert.True(t, ok)
	assert.Equal(t, []string{"1"}, value)
}
