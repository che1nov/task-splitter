package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

// slogLogger реализация Logger с использованием slog
type slogLogger struct {
	logger *slog.Logger
}

// New создает новый логгер
func New(level string) Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &slogLogger{logger: logger}
}

// Debug логирует сообщение уровня Debug
func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// DebugContext логирует сообщение уровня Debug с контекстом
func (l *slogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Info логирует сообщение уровня Info
func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// InfoContext логирует сообщение уровня Info с контекстом
func (l *slogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Warn логирует сообщение уровня Warn
func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// WarnContext логирует сообщение уровня Warn с контекстом
func (l *slogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Error логирует сообщение уровня Error
func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// ErrorContext логирует сообщение уровня Error с контекстом
func (l *slogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// With создает новый логгер с дополнительными полями
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{logger: l.logger.With(args...)}
}

