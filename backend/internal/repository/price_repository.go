package repository

import (
	"context"

	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type PriceRepository interface {
	GetLatestPrices(ctx context.Context) (types.Prices, error)
	Save(ctx context.Context, prices types.Prices) error
}

type priceRepo struct {
	db database.DB
}

func NewPriceRepository(db database.DB) PriceRepository {
	return &priceRepo{db: db}
}

func (r *priceRepo) GetLatestPrices(ctx context.Context) (types.Prices, error) {
	const query = `
	SELECT bs_amount, usd_amount
	FROM prices
	ORDER BY created_at DESC
	LIMIT 1
	`
	var price types.Prices
	err := r.db.QueryRow(ctx, query).Scan(&price.BsAmount, &price.UsdAmount)
	return price, err
}

func (r *priceRepo) Save(ctx context.Context, prices types.Prices) error {
	const query = `
	INSERT INTO prices (bs_amount, usd_amount)
	VALUES ($1, $2)
	`
	return r.db.ExecContext(ctx, query, prices.BsAmount, prices.UsdAmount)
}
