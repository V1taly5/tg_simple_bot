package sl

import (
	"context"
	"log/slog"
)

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value("logger").(*slog.Logger)
	if !ok {
		slog.Info("logger not found")
		return nil
	}
	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
