package repository

import (
	"context"

	"rifa/backend/pkg/db"
)

type IdempotencyRecord struct {
	Key          string
	BodyHash     string
	StatusCode   int
	ResponseBody []byte
}

type IdempotencyRepository interface {
	Get(ctx context.Context, key string) (IdempotencyRecord, error)
	Save(ctx context.Context, rec IdempotencyRecord) error
}

type idempotencyRepo struct{ db db.DB }

func NewIdempotencyRepository(db db.DB) IdempotencyRepository {
	return &idempotencyRepo{db: db}
}

func (r *idempotencyRepo) Get(
	ctx context.Context,
	key string,
) (IdempotencyRecord, error) {
	row := r.db.QueryRow(ctx, `
		SELECT key, body_hash, status_code, response_body
		FROM idempotency_keys
		WHERE key = $1
	`, key)

	var rec IdempotencyRecord
	err := row.Scan(&rec.Key, &rec.BodyHash, &rec.StatusCode, &rec.ResponseBody)
	if err != nil {
		return IdempotencyRecord{}, err
	}

	return rec, nil
}

func (r *idempotencyRepo) Save(ctx context.Context, rec IdempotencyRecord) error {
	err := r.db.ExecContext(ctx, `
		INSERT INTO idempotency_keys (key, body_hash, status_code, response_body)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO NOTHING
	`, rec.Key, rec.BodyHash, rec.StatusCode, rec.ResponseBody)
	if err != nil {
		return err
	}
	return nil
}
