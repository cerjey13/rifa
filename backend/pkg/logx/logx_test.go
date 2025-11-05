//go:build !integration
// +build !integration

package logx

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestSlogLoggerLevels(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		logFunc  func(Logger)
		expected string
	}{
		{
			name: "debug in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Debug(context.Background(), "debug message", "key", "value")
			},
			expected: "debug message",
		},
		{
			name: "info in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Info(context.Background(), "info message", "key", "value")
			},
			expected: "info message",
		},
		{
			name: "warn in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Warn(context.Background(), "warn message", "key", "value")
			},
			expected: "warn message",
		},
		{
			name: "error in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Error(context.Background(), "error message", "key", "value")
			},
			expected: "error message",
		},
		{
			name: "json format in production",
			env:  "production",
			logFunc: func(l Logger) {
				l.Info(context.Background(), "json message", "key", "value")
			},
			expected: `"msg":"json message"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var handler slog.Handler

			if tt.env == "production" {
				handler = slog.NewJSONHandler(&buf, nil)
			} else {
				handler = slog.NewTextHandler(
					&buf,
					&slog.HandlerOptions{Level: slog.LevelDebug},
				)
			}

			logger := &logger{base: slog.New(handler)}

			tt.logFunc(logger)

			got := buf.String()
			if !strings.Contains(got, tt.expected) {
				t.Errorf(
					"expected log output to contain %q, got %q",
					tt.expected,
					got,
				)
			}
		})
	}
}
