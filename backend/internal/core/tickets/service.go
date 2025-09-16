package ticket

import (
	"context"

	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
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
	repo repository.TicketRepository
}

func NewService(db database.DB) Service {
	return &service{
		repo: repository.NewTicketRepository(db),
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
	return s.repo.AssignTickets(
		ctx,
		lotteryID,
		userID,
		purchaseID,
		selectedNumbers,
		quantity,
	)
}

func (s *service) SearchTickets(
	ctx context.Context,
	tickets []int,
) ([]int, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUnavailableNumbers(ctx, lotteryID, tickets)
}

func (s *service) GetAvailability(ctx context.Context) (float64, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		return 0, err
	}
	return s.repo.GetAvailabilityPercentage(ctx, lotteryID)
}

func (s *service) GetUserTickets(
	ctx context.Context,
	userID string,
) ([]int, error) {
	lotteryID, err := s.repo.GetActiveLotteryID(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserTickets(ctx, userID, lotteryID)
}
