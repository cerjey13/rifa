package logx

import (
	"bytes"
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
				l.Debug("debug message", "key", "value")
			},
			expected: "debug message",
		},
		{
			name: "info in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Info("info message", "key", "value")
			},
			expected: "info message",
		},
		{
			name: "warn in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Warn("warn message", "key", "value")
			},
			expected: "warn message",
		},
		{
			name: "error in dev",
			env:  "dev",
			logFunc: func(l Logger) {
				l.Error("error message", "key", "value")
			},
			expected: "error message",
		},
		{
			name: "json format in production",
			env:  "production",
			logFunc: func(l Logger) {
				l.Info("json message", "key", "value")
			},
			expected: `"msg":"json message"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var handler slog.Handler

			// replace os.Stdout with our buffer for testing
			if tt.env == "production" {
				handler = slog.NewJSONHandler(&buf, nil)
			} else {
				handler = slog.NewTextHandler(
					&buf,
					&slog.HandlerOptions{Level: slog.LevelDebug},
				)
			}

			logger := &slogLogger{l: slog.New(handler)}

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
