package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewLogger(env string) corelog.Logger {
	var handler slog.Handler
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	if strings.EqualFold(env, "prod") || strings.EqualFold(env, "production") {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	return &SlogLogger{logger: slog.New(handler)}
}

func (l *SlogLogger) With(attrs ...any) corelog.Logger {
	return &SlogLogger{logger: l.logger.With(attrs...)}
}

func (l *SlogLogger) Info(ctx context.Context, msg string, attrs ...any) {
	l.logger.InfoContext(ctx, msg, attrs...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, attrs ...any) {
	l.logger.ErrorContext(ctx, msg, attrs...)
}
