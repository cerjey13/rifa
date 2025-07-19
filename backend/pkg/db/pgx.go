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

type pgxRowsAdapter struct {
	rows pgx.Rows
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

func (d *pgxpoolAdapter) Query(
	ctx context.Context,
	query string,
	args ...any,
) (Rows, error) {
	rows, err := d.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRowsAdapter{rows: rows}, nil
}
func (d *pgxpoolAdapter) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := d.pool.Exec(ctx, query, args...)
	return err
}

func (r *pgxRowsAdapter) Next() bool             { return r.rows.Next() }
func (r *pgxRowsAdapter) Scan(dest ...any) error { return r.rows.Scan(dest...) }
func (r *pgxRowsAdapter) Close()                 { r.rows.Close() }
func (r *pgxRowsAdapter) Err() error             { return r.rows.Err() }
