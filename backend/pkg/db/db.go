package db

import (
	"context"
	"fmt"
	"path/filepath"
	"rifa/backend/pkg/config"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	ExecContext(ctx context.Context, query string, args ...any) error
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

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		pool, err = pgxpool.New(ctx, cfg.DatabaseUrl)
		if err != nil {
			return
		}
		err = pool.Ping(ctx)
		if err != nil {
			return
		}
		err = runMigrations(cfg)
	})

	return pool, err
}

func runMigrations(cfg *config.Config) error {
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
	defer m.Close()
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
