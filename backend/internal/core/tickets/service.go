package ticket

import (
	"context"

	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
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
}

type service struct {
	repo repository.TicketRepository
}

func NewService(repo repository.TicketRepository) Service {
	return &service{repo: repo}
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
