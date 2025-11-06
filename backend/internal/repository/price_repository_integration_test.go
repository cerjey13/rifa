//go:build integration

package repository

import (
	"context"
	"math"
	"testing"
	"time"

	"rifa/backend/internal/types"
	"rifa/backend/testutils"
)

func TestPriceRepository_CreateAndGetLatest(t *testing.T) {
	pg := testutils.StartPostgres(t)
	defer testutils.StopPostgres(t, pg)

	repo := NewPriceRepository(pg.Database)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		input          types.Prices
		expectedLatest types.Prices
	}{
		{
			name:  "Initial default record should exist from migration",
			input: types.Prices{},
			expectedLatest: types.Prices{
				BsAmount:  150,
				UsdAmount: 1,
			},
		},
		{
			name: "Insert new price overrides previous one",
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
			name: "Second insert should become latest",
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
				err := repo.Save(ctx, tt.input)
				if err != nil {
					t.Fatalf("save should succeed, got error: %v", err)
				}
				time.Sleep(2 * time.Millisecond)
			}

			got, err := repo.GetLatestPrices(ctx)
			if err != nil {
				t.Fatalf("expected no error fetching latest, got: %v", err)
			}

			if diff := math.Abs(tt.expectedLatest.BsAmount - got.BsAmount); diff > 0.0001 {
				t.Errorf("BsAmount mismatch: expected %.4f, got %.4f (diff=%.4f)",
					tt.expectedLatest.BsAmount, got.BsAmount, diff)
			}

			if diff := math.Abs(tt.expectedLatest.UsdAmount - got.UsdAmount); diff > 0.0001 {
				t.Errorf("UsdAmount mismatch: expected %.4f, got %.4f (diff=%.4f)",
					tt.expectedLatest.UsdAmount, got.UsdAmount, diff)
			}
		})
	}
}
