package sl

import (
	"log/slog"
)

type Logger interface {
	Info(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	Debug(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
	Fatal(msg string, attrs ...slog.Attr)
	With(attrs ...slog.Attr) Logger
}

type MyLogger struct {
	*slog.Logger
	fatalFunc func(msg string, attrs ...any)
}

func New(logger *slog.Logger, fatalFunc func(msg string, attrs ...any)) *MyLogger {
	return &MyLogger{Logger: logger, fatalFunc: fatalFunc}
}

func (l *MyLogger) Info(msg string, attrs ...slog.Attr) {
	l.Logger.Info(msg, attrsToAny(attrs)...)
}

func (l *MyLogger) Error(msg string, attrs ...slog.Attr) {
	l.Logger.Error(msg, attrsToAny(attrs)...)
}

func (l *MyLogger) Debug(msg string, attrs ...slog.Attr) {
	l.Logger.Debug(msg, attrsToAny(attrs)...)
}

func (l *MyLogger) Warn(msg string, attrs ...slog.Attr) {
	l.Logger.Warn(msg, attrsToAny(attrs)...)
}

func (l *MyLogger) Fatal(msg string, attrs ...slog.Attr) {
	l.fatalFunc(msg, attrsToAny(attrs)...)
}

func (l *MyLogger) With(attrs ...slog.Attr) Logger {
	return &MyLogger{Logger: l.Logger.With(attrsToAny(attrs)...)}
}

func attrsToAny(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return args
}
