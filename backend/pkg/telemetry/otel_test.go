package telemetry

import (
	"context"
	"testing"

	"rifa/backend/pkg/config"
)

func TestInitCollector(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		opts     config.CollectorOpts
		wantErr  bool
		wantNoop bool
	}{
		{
			name: "development mode skips collector",
			opts: config.CollectorOpts{
				CollectorEnv: "development",
			},
			wantErr:  false,
			wantNoop: true,
		},
		{
			name: "invalid endpoint triggers error",
			opts: config.CollectorOpts{
				CollectorEnv:      "production",
				CollectorExporter: "invalid:://endpoint",
			},
			wantErr:  true,
			wantNoop: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown, err := InitCollector(ctx, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitCollector() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("InitCollector() unexpected error: %v", err)
			}

			if shutdown == nil {
				t.Fatalf("InitCollector() returned nil shutdown func")
			}

			// Call shutdown to ensure it's safe (no panic, no error)
			if err := shutdown(ctx); err != nil {
				t.Errorf("shutdown returned error: %v", err)
			}

			if tt.wantNoop {
				// In dev mode the shutdown is a simple no-op
				// We can't directly test internal providers,
				// but we ensure it executes without error.
			}
		})
	}
}

func TestParseKeyValueList(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]string
	}{
		{
			name: "empty string",
			in:   "",
			want: map[string]string{},
		},
		{
			name: "single key=value",
			in:   "foo=bar",
			want: map[string]string{"foo": "bar"},
		},
		{
			name: "multiple pairs with spaces",
			in:   "foo=bar, baz=qux",
			want: map[string]string{"foo": "bar", "baz": "qux"},
		},
		{
			name: "skip malformed pairs",
			in:   "good=pair,badpair,nokey=,=noval",
			want: map[string]string{"good": "pair"},
		},
		{
			name: "trim spaces around keys and values",
			in:   " a = b , c = d ",
			want: map[string]string{"a": "b", "c": "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseKeyValueList(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d pairs, got %d", len(tt.want), len(got))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("key %q: expected %q, got %q", k, v, got[k])
				}
			}
		})
	}
}
