package db

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"rifa/backend/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	ExecContext(ctx context.Context, query string, args ...any) error
	BeginTx(ctx context.Context) (Tx, error)
	Close()
}

type Tx interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	ExecContext(ctx context.Context, query string, args ...any) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type Driver interface {
	Open(ctx context.Context, dsn string) (DB, error)
}

var (
	db      DB
	initErr error
	once    sync.Once
)

func Connect(
	ctx context.Context,
	drv Driver,
	cfg *config.DatabaseOpts,
) (DB, error) {
	once.Do(func() {
		var err error
		db, err = drv.Open(ctx, cfg.DatabaseUrl)
		if err != nil {
			initErr = err
			return
		}

		err = runMigrations(cfg)
		if err != nil {
			initErr = err
			return
		}
	})

	return db, initErr
}

func runMigrations(cfg *config.DatabaseOpts) error {
	absMigrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return err
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", absMigrationsPath),
		cfg.DatabaseUrl,
	)
	if err != nil {
		return err
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migration source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migration db close error: %v", dbErr)
		}
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
