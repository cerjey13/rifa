package ticket

import (
	"context"

	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
)

type Service interface {
	BuyTickets(
		ctx context.Context,
		lotteryID,
		userID,
		purchaseID string,
		selectedNumbers []string,
		quantity int,
	) ([]types.Ticket, error)
	SearchTickets(ctx context.Context, tickets []int) ([]int, error)
	GetAvailability(ctx context.Context) (float64, error)
	GetUserTickets(ctx context.Context, userID string) ([]int, error)
}

type service struct {
	repo   repository.TicketRepository
	logger logx.Logger
}

func NewService(db database.DB, logger logx.Logger) Service {
	return &service{
		repo:   repository.NewTicketRepository(db),
		logger: logger,
	}
}

func (s *service) BuyTickets(
	ctx context.Context,
	lotteryID,
	userID,
	purchaseID string,
	selectedNumbers []string,
	quantity int,
) ([]types.Ticket, error) {
	tickets, err := s.repo.AssignTickets(
		ctx,
		lotteryID,
		userID,
		purchaseID,
		selectedNumbers,
		quantity,
	)
	if err != nil {
		s.logger.Error(
			"Failed to purchase tickets",
			"user",
			userID,
			"error",
			err,
		)
		return nil, err
	}

	return tickets, nil
}

func (s *service) SearchTickets(
	ctx context.Context,
	tickets []int,
) ([]int, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		s.logger.Error("Failed to get active lottery", "error", err)
		return nil, err
	}

	numbers, err := s.repo.GetUnavailableNumbers(ctx, lotteryID, tickets)
	if err != nil {
		s.logger.Error(
			"Failed to check if the selected numbers are unavailable",
			"error",
			err,
		)
		return nil, err
	}

	return numbers, nil
}

func (s *service) GetAvailability(ctx context.Context) (float64, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		s.logger.Error("Failed to get active lottery", "error", err)
		return 0, err
	}

	percentage, err := s.repo.GetAvailabilityPercentage(ctx, lotteryID)
	if err != nil {
		s.logger.Error("Failed to get")
		return 0, err
	}

	return percentage, nil
}

func (s *service) GetUserTickets(
	ctx context.Context,
	userID string,
) ([]int, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		s.logger.Error("Failed to get active lottery", "error", err)
		return nil, err
	}

	userTickets, err := s.repo.GetUserTickets(ctx, userID, lotteryID)
	if err != nil {
		s.logger.Error(
			"Failed to retrieve user buyed tickets",
			"user",
			userID,
			"error",
			err,
		)
		return nil, err
	}

	return userTickets, nil
}
