//go:build integration

package testutils

import (
	"context"
	"testing"
)

func TestStartAndStopPostgres(t *testing.T) {
	st := StartPostgres(t)
	defer StopPostgres(t, st)

	if err := st.Database.Ping(context.Background()); err != nil {
		t.Fatalf("database should be reachable: %v", err)
	}

	rows, err := st.Database.Query(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("expected at least one row from SELECT 1")
	}

	var val int
	if err := rows.Scan(&val); err != nil {
		t.Fatalf("failed to scan result: %v", err)
	}

	if val != 1 {
		t.Errorf("expected SELECT 1 to return 1, got %d", val)
	}
}
