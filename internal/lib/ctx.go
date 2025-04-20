package lib

import (
	"context"
	"errors"
)

var ErrContextValueMissing = errors.New("context value missing or wrong type")

func GetFromContext[T any](ctx context.Context, key any) (T, error) {
	val := ctx.Value(key)
	if val == nil {
		var zero T
		return zero, ErrContextValueMissing
	}
	v, ok := val.(T)
	if !ok {
		var zero T
		return zero, ErrContextValueMissing
	}
	return v, nil
}
