package logx

import (
	"log/slog"
	"os"
	"runtime"
	"strconv"
)

// Logger is a minimal interface for structured logging.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type slogLogger struct {
	l *slog.Logger
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }

func NewLogger(env string) Logger {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					if _, file, line, ok := runtime.Caller(7); ok {
						a.Value = slog.StringValue(
							file + ":" + strconv.Itoa(line),
						)
					}
				}
				return a
			},
		})
	}

	return &slogLogger{l: slog.New(handler)}
}
