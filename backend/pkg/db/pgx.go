package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxpoolAdapter struct {
	pool *pgxpool.Pool
}

type pgxRowAdapter struct {
	row pgx.Row
}

func (r *pgxRowAdapter) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func NewPgxpoolAdapter(pool *pgxpool.Pool) DB {
	return &pgxpoolAdapter{pool: pool}
}

func (d *pgxpoolAdapter) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) Row {
	return &pgxRowAdapter{row: d.pool.QueryRow(ctx, query, args...)}
}

func (d *pgxpoolAdapter) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := d.pool.Exec(ctx, query, args...)
	return err
}
