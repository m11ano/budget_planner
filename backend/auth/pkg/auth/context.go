package auth

import (
	"context"

	"github.com/LastPossum/kamino"
)

type contextKey string

const (
	ContextKeyAuthData       contextKey = "auth_context_auth_data"
	ContextKeyAuthCheckRight contextKey = "auth_context_auth_check_right"
)

func GetAuthData(ctx context.Context) *AuthData {
	data, ok := ctx.Value(ContextKeyAuthData).(*AuthData)
	if !ok {
		return nil
	}

	copy, err := kamino.Clone(data)
	if err != nil {
		panic(err)
	}

	return copy
}

func SetAuthData(ctx context.Context, data *AuthData) context.Context {
	return context.WithValue(ctx, ContextKeyAuthData, data)
}

func WithoutCheckRight(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyAuthCheckRight, false)
}

func WithCheckRight(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyAuthCheckRight, true)
}

func IsNeedToCheckRights(ctx context.Context) bool {
	is, ok := ctx.Value(ContextKeyAuthCheckRight).(bool)
	return ok && is
}
