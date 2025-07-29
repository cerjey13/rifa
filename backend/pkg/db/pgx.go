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

type pgxTxAdapter struct {
	tx pgx.Tx
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

func (d *pgxpoolAdapter) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTxAdapter{tx: tx}, nil
}

func (r *pgxRowsAdapter) Next() bool             { return r.rows.Next() }
func (r *pgxRowsAdapter) Scan(dest ...any) error { return r.rows.Scan(dest...) }
func (r *pgxRowsAdapter) Close()                 { r.rows.Close() }
func (r *pgxRowsAdapter) Err() error             { return r.rows.Err() }

func (t *pgxTxAdapter) Query(
	ctx context.Context,
	query string,
	args ...any,
) (Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRowsAdapter{rows: rows}, nil
}

func (t *pgxTxAdapter) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) Row {
	return &pgxRowAdapter{row: t.tx.QueryRow(ctx, query, args...)}
}

func (t *pgxTxAdapter) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := t.tx.Exec(ctx, query, args...)
	return err
}

func (t *pgxTxAdapter) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTxAdapter) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
