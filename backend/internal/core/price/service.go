package price

import (
	"context"

	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type Service interface {
	GetPrices(ctx context.Context) (types.Prices, error)
	Update(ctx context.Context, bs, usd float64) error
}

type service struct {
	repo repository.PriceRepository
}

func NewService(db database.DB) Service {
	return &service{
		repo: repository.NewPriceRepository(db),
	}
}

func (s *service) GetPrices(ctx context.Context) (types.Prices, error) {
	return s.repo.GetLatestPrices(ctx)
}

func (s *service) Update(ctx context.Context, bs, usd float64) error {
	return s.repo.Save(ctx, types.Prices{BsAmount: bs, UsdAmount: usd})
}
