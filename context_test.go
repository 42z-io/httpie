package httpie

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testContextStruct struct {
	Name string
}

type ctxKeyType int

var uniqueCtxKey ctxKeyType = 0
var otherCtxKey ctxKeyType = 1

func TestGetContextValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), uniqueCtxKey, &testContextStruct{Name: "test"})
	result := GetContextValue[testContextStruct](ctx, uniqueCtxKey)
	assert.NotNil(t, result)
	assert.Equal(t, "test", result.Name)
}

func TestGetContextValueNotExists(t *testing.T) {
	ctx := context.WithValue(context.Background(), uniqueCtxKey, &testContextStruct{Name: "test"})
	result := GetContextValue[testContextStruct](ctx, otherCtxKey)
	assert.Nil(t, result)
}

func TestGetContextValueCastError(t *testing.T) {
	ctx := context.WithValue(context.Background(), uniqueCtxKey, &testContextStruct{Name: "test"})
	result := GetContextValue[map[string]string](ctx, uniqueCtxKey)
	assert.Nil(t, result)
}
