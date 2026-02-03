package log

import "context"

type Logger interface {
	With(attrs ...any) Logger
	Info(ctx context.Context, msg string, attrs ...any)
	Error(ctx context.Context, msg string, attrs ...any)
}

type noopLogger struct{}

func (noopLogger) With(_ ...any) Logger { return noopLogger{} }
func (noopLogger) Info(_ context.Context, _ string, _ ...any) {
}
func (noopLogger) Error(_ context.Context, _ string, _ ...any) {
}

type loggerKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		logger = noopLogger{}
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return logger
	}
	return noopLogger{}
}

func FromContextOr(ctx context.Context, fallback Logger) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return logger
	}
	if fallback == nil {
		return noopLogger{}
	}
	return fallback
}
