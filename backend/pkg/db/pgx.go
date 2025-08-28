package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDriver struct{}

func NewPostgresDriver() Driver { return &postgresDriver{} }

func (postgresDriver) Open(ctx context.Context, dsn string) (DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	// Wrap pgxpool with your adapter so callers use your DB interface
	return NewPGX(pool), nil
}

type PGXPool struct {
	Pool *pgxpool.Pool
}

func NewPGX(pool *pgxpool.Pool) *PGXPool { return &PGXPool{Pool: pool} }

func (p *PGXPool) Query(ctx context.Context, q string, args ...any) (Rows, error) {
	return p.Pool.Query(ctx, q, args...)
}

func (p *PGXPool) QueryRow(ctx context.Context, q string, args ...any) Row {
	return p.Pool.QueryRow(ctx, q, args...)
}

func (p *PGXPool) ExecContext(ctx context.Context, q string, args ...any) error {
	_, err := p.Pool.Exec(ctx, q, args...)
	return err
}

func (p *PGXPool) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTx{tx: tx}, nil
}

func (p *PGXPool) Close() {
	p.Pool.Close()
}

type pgxTx struct{ tx pgx.Tx }

func (t *pgxTx) Query(ctx context.Context, q string, args ...any) (Rows, error) {
	return t.tx.Query(ctx, q, args...)
}

func (t *pgxTx) QueryRow(ctx context.Context, q string, args ...any) Row {
	return t.tx.QueryRow(ctx, q, args...)
}

func (t *pgxTx) ExecContext(ctx context.Context, q string, args ...any) error {
	_, err := t.tx.Exec(ctx, q, args...)
	return err
}

func (t *pgxTx) Commit(ctx context.Context) error   { return t.tx.Commit(ctx) }
func (t *pgxTx) Rollback(ctx context.Context) error { return t.tx.Rollback(ctx) }
