package logx

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"

	"go.opentelemetry.io/otel/trace"
)

// Logger is a minimal interface for structured logging.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type logger struct {
	base *slog.Logger
}

func (s *logger) Debug(ctx context.Context, msg string, args ...any) {
	s.base.Debug(msg, enrichWithTrace(ctx, args)...)
}
func (s *logger) Info(ctx context.Context, msg string, args ...any) {
	s.base.Info(msg, enrichWithTrace(ctx, args)...)
}
func (s *logger) Warn(ctx context.Context, msg string, args ...any) {
	s.base.Warn(msg, enrichWithTrace(ctx, args)...)
}
func (s *logger) Error(ctx context.Context, msg string, args ...any) {
	s.base.Error(msg, enrichWithTrace(ctx, args)...)
}

func enrichWithTrace(ctx context.Context, args []any) []any {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		args = append(args,
			"trace_id", spanCtx.TraceID().String(),
			"span_id", spanCtx.SpanID().String(),
		)
	}
	return args
}

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

	return &logger{base: slog.New(handler)}
}
