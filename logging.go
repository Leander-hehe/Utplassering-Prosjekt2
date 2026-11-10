package main

import (
	"context"
	"log/slog"
)

type ContextKey string

const LoggerKey ContextKey = "logger"

func (app *application) LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)

	if !ok {
		return slog.Default()
	}
	return logger
}

func (app *application) ContextLogger(ctx context.Context, attr slog.Attr) (context.Context, *slog.Logger) {
	logger := app.LoggerFromContext(ctx).With(attr)
	return app.WithLogger(ctx, logger), logger
}

func (app *application) WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
