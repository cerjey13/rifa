//go:build integration

package price

import (
	"context"
	"math"
	"testing"
	"time"

	"rifa/backend/internal/types"
	"rifa/backend/pkg/logx"
	"rifa/backend/testutils"
)

func TestPriceService_Integration(t *testing.T) {
	pg := testutils.StartPostgres(t)
	defer testutils.StopPostgres(t, pg)

	logger := logx.NewLogger("development")
	svc := NewService(pg.Database, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		input          types.Prices
		expectError    bool
		expectedLatest types.Prices
	}{
		{
			name: "Initial default price exists from migration",
			expectedLatest: types.Prices{
				BsAmount:  150,
				UsdAmount: 1,
			},
		},
		{
			name: "Update should insert first custom price",
			input: types.Prices{
				BsAmount:  37.5,
				UsdAmount: 38.2,
			},
			expectedLatest: types.Prices{
				BsAmount:  37.5,
				UsdAmount: 38.2,
			},
		},
		{
			name: "Second update should override latest price",
			input: types.Prices{
				BsAmount:  38.7,
				UsdAmount: 39.1,
			},
			expectedLatest: types.Prices{
				BsAmount:  38.7,
				UsdAmount: 39.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.BsAmount != 0 && tt.input.UsdAmount != 0 {
				err := svc.Update(ctx, tt.input.BsAmount, tt.input.UsdAmount)
				if err != nil {
					t.Fatalf("update failed: %v", err)
				}
				time.Sleep(2 * time.Millisecond)
			}

			got, err := svc.GetPrices(ctx)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error fetching latest: %v", err)
			}

			bsDiff := math.Abs(tt.expectedLatest.BsAmount - got.BsAmount)
			if bsDiff > 0.0001 {
				t.Errorf(
					"BsAmount mismatch: expected %.4f, got %.4f (diff=%.4f)",
					tt.expectedLatest.BsAmount,
					got.BsAmount,
					bsDiff,
				)
			}

			usDiff := math.Abs(tt.expectedLatest.UsdAmount - got.UsdAmount)
			if usDiff > 0.0001 {
				t.Errorf(
					"UsdAmount mismatch: expected %.4f, got %.4f (diff=%.4f)",
					tt.expectedLatest.UsdAmount,
					got.UsdAmount,
					usDiff,
				)
			}
		})
	}
}
