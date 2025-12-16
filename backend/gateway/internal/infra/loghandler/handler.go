package loghandler

import (
	"context"
	"log/slog"
	"runtime"
)

type handlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *handlerMiddlware {
	return &handlerMiddlware{next: next}
}

func (h *handlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *handlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	data, ok := GetData(ctx)
	if ok {
		for k, v := range data {
			rec.Add(k, v)
		}
	}

	if IsWithSource(ctx) {
		if pc, file, line, ok := runtime.Caller(3); ok {
			fn := runtime.FuncForPC(pc).Name()

			rec.AddAttrs(
				slog.Group("source",
					slog.String("file", file),
					slog.String("function", fn),
					slog.Int("line", line),
				),
			)
		}
	}

	return h.next.Handle(ctx, rec)
}

func (h *handlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handlerMiddlware{next: h.next.WithAttrs(attrs)}
}

func (h *handlerMiddlware) WithGroup(name string) slog.Handler {
	return &handlerMiddlware{next: h.next.WithGroup(name)}
}
