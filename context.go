package httpie

import "context"

// GetContextValue returns the value associated with this context for key, or nil if no value is associated with key.
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
