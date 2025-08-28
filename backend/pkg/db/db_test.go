package db

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"rifa/backend/pkg/config"
)

type fakeDB struct {
	closed int
}

func (f *fakeDB) Query(context.Context, string, ...any) (Rows, error) { return nil, nil }
func (f *fakeDB) QueryRow(context.Context, string, ...any) Row        { return nil }
func (f *fakeDB) ExecContext(context.Context, string, ...any) error   { return nil }
func (f *fakeDB) BeginTx(context.Context) (Tx, error)                 { return nil, nil }
func (f *fakeDB) Close()                                              { f.closed++ }

type fakeDriver struct {
	openCalls  int
	dbToReturn DB
	err        error
}

func (d *fakeDriver) Open(ctx context.Context, dsn string) (DB, error) {
	d.openCalls++
	if d.err != nil {
		return nil, d.err
	}
	return d.dbToReturn, nil
}

func resetSingleton() {
	db = nil
	initErr = nil
	once = sync.Once{}
}

// If the driver fails, Connect should return (nil, err) and cache the error.
// Subsequent calls must NOT re-open and must return the same error.
func TestConnect_DriverError_IsCached(t *testing.T) {
	resetSingleton()
	t.Cleanup(resetSingleton)

	drv := &fakeDriver{err: errors.New("open failed")}
	cfg := &config.Config{DatabaseUrl: "postgres://does-not-matter"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 1st call fails
	db1, err1 := Connect(ctx, drv, cfg)
	if err1 == nil {
		t.Fatalf("expected error on first Connect, got nil")
	}
	if db1 != nil {
		t.Fatalf("expected nil DB on error, got %T", db1)
	}
	if drv.openCalls != 1 {
		t.Fatalf("driver.Open called %d times, want 1", drv.openCalls)
	}

	// 2nd call should not try to open again; same cached error
	db2, err2 := Connect(ctx, drv, cfg)
	if err2 == nil {
		t.Fatalf("expected cached error on second Connect, got nil")
	}
	if db2 != nil {
		t.Fatalf("expected nil DB on error, got %T", db2)
	}
	if drv.openCalls != 1 {
		t.Fatalf("driver.Open called %d times (should be once)", drv.openCalls)
	}
}

// If migrations fail after a successful open, Connect should return an error.
// With your current code, the DB remains set; later calls return the same DB and the same error.
// Also ensure driver.Open was called only once.
func TestConnect_MigrationError_ReturnsError_AndSingleton(t *testing.T) {
	resetSingleton()
	t.Cleanup(resetSingleton)

	fdb := &fakeDB{}
	drv := &fakeDriver{dbToReturn: fdb}

	cfg := &config.Config{
		DatabaseUrl: "postgres://user:pass@localhost:5432/db?sslmode=disable",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1st call: driver open succeeds, migrations fail -> error returned.
	gotDB1, err1 := Connect(ctx, drv, cfg)
	if err1 == nil {
		t.Fatalf("expected migration error, got nil")
	}

	if gotDB1 == nil {
		t.Fatalf("expected non-nil DB instance when open succeeded")
	}
	if drv.openCalls != 1 {
		t.Fatalf("driver.Open calls = %d, want 1", drv.openCalls)
	}
	// The returned DB should be our fake instance.
	if got, ok := gotDB1.(*fakeDB); !ok || got != fdb {
		t.Fatalf("returned DB is not the same fakeDB instance")
	}

	if fdb.closed != 0 {
		t.Fatalf("unexpected Close() call on migration error")
	}

	// 2nd call: once.Do prevents another open; same DB & same error are returned.
	gotDB2, err2 := Connect(ctx, drv, cfg)
	if err2 == nil {
		t.Fatalf("expected cached migration error on second Connect, got nil")
	}
	if gotDB2 == nil {
		t.Fatalf("expected non-nil DB on second call")
	}
	if drv.openCalls != 1 {
		t.Fatalf("driver.Open called %d times, want 1", drv.openCalls)
	}
	if got2, ok := gotDB2.(*fakeDB); !ok || got2 != fdb {
		t.Fatalf("second returned DB not the same instance")
	}
}

func TestConnect_Singleton_Concurrency(t *testing.T) {
	resetSingleton()
	t.Cleanup(resetSingleton)

	fdb := &fakeDB{}
	drv := &fakeDriver{dbToReturn: fdb}
	cfg := &config.Config{
		DatabaseUrl: "postgres://user:pass@localhost:5432/db?sslmode=disable",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)

	errs := make(chan error, N)
	dbs := make(chan DB, N)

	for _ = range N {
		go func() {
			defer wg.Done()
			d, err := Connect(ctx, drv, cfg)
			errs <- err
			dbs <- d
		}()
	}
	wg.Wait()
	close(errs)
	close(dbs)

	// Exactly one Open call due to sync.Once
	if drv.openCalls != 1 {
		t.Fatalf("driver.Open called %d times, want 1", drv.openCalls)
	}

	// All calls should have the same outcome (migration error) and same DB pointer.
	for e := range errs {
		if e == nil {
			t.Fatalf("expected error (migration failure), got nil")
		}
	}
	for d := range dbs {
		if d == nil {
			t.Fatalf("expected non-nil DB instance")
		}
		if got, ok := d.(*fakeDB); !ok || got != fdb {
			t.Fatalf("returned DB not the same singleton instance")
		}
	}
}
