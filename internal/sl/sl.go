package sl

import (
	"log/slog"
)

type LoggerFunc  func(msg string, attrs ...slog.Attr)

type MyLogger struct {
	*slog.Logger
	Fatal  LoggerFunc
}

func New(logger *slog.Logger, fatalFunc LoggerFunc) *MyLogger {
	return &MyLogger{
		Logger: logger,
		Fatal:  fatalFunc,
	}
}
