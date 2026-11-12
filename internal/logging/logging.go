package logging

import (
	"context"
	"log/slog"
)

type ContextKey string

const LoggerKey ContextKey = "logger"

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)

	if !ok {
		return slog.Default()
	}
	return logger
}

func ContextLogger(ctx context.Context, attr slog.Attr) (context.Context, *slog.Logger) {
	logger := LoggerFromContext(ctx).With(attr)
	return WithLogger(ctx, logger), logger
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
