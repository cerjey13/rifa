package price

import (
	"context"
	"errors"

	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
)

type Service interface {
	GetPrices(ctx context.Context) (types.Prices, error)
	Update(ctx context.Context, bs, usd float64) error
}

type service struct {
	repo   repository.PriceRepository
	logger logx.Logger
}

func NewService(db database.DB, logger logx.Logger) Service {
	return &service{
		repo:   repository.NewPriceRepository(db),
		logger: logger,
	}
}

func (s *service) GetPrices(ctx context.Context) (types.Prices, error) {
	prices, err := s.repo.GetLatestPrices(ctx)
	if err != nil {
		s.logger.Error("Failed to get latest prices", "error", err)
		return types.Prices{}, err
	}
	return prices, nil
}

func (s *service) Update(ctx context.Context, bs, usd float64) error {
	if bs <= 0 || usd <= 0 {
		s.logger.Warn(
			"Attempt to update prices with non-positive values",
			"bs",
			bs,
			"usd",
			usd,
		)
		return errors.New("valores inválidos")
	}

	err := s.repo.Save(ctx, types.Prices{BsAmount: bs, UsdAmount: usd})
	if err != nil {
		s.logger.Error("Failed to save prices data into the db", "error", err)
		return err
	}

	s.logger.Info("Updated prices successfully", "bs", bs, "usd", usd)
	return nil
}
