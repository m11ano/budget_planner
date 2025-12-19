package loghandler

import "context"

type contextKey string

const (
	contextKeyData       contextKey = "log_context_data"
	contextKeyWithSource contextKey = "log_context_with_source"
)

type contextData map[string]any

func GetData(ctx context.Context) (contextData, bool) {
	data, ok := ctx.Value(contextKeyData).(contextData)
	return data, ok
}

func SetData(ctx context.Context, key string, value any) (contextData, contextKey) {
	data, ok := GetData(ctx)
	if !ok {
		data = contextData{}
	}

	data[key] = value

	return data, contextKeyData
}

func SetContextData(ctx context.Context, key string, value any) context.Context {
	data, ctxKey := SetData(ctx, key, value)
	return context.WithValue(ctx, ctxKey, data)
}

func WithSource(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyWithSource, true)
}

func IsWithSource(ctx context.Context) bool {
	is, ok := ctx.Value(contextKeyWithSource).(bool)
	return ok && is
}
