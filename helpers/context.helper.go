package helpers

import (
	"context"
	"errors"
)

type ctxKey string

func NewctxKey(key string) ctxKey {
	return ctxKey(key)
}

func SetInContext(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, NewctxKey(key), value)
}

func GetFromContext(ctx context.Context, key string) (interface{}, error) {
	val := ctx.Value(NewctxKey(key))
	if val == nil {
		return nil, errors.New("value not found in context")
	}
	return val, nil
}
