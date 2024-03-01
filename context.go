package httpie

import "context"

func GetContextValue[T any](ctx context.Context, key any) *T {
	value := ctx.Value(key)
	if value == nil {
		return nil
	}
	result, ok := value.(*T)
	if !ok {
		return nil
	}
	return result
}
