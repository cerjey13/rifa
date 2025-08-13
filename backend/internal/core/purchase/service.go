package purchase

import (
	"context"
	"fmt"
	"log"
	"time"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/core/email"
	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"
)

type Service interface {
	Create(ctx context.Context, req *form.CreatePurchaseRequest) error
	GetAll(
		ctx context.Context,
		filters dto.GetAllPurchases,
	) ([]form.Purchases, int, error)
	UpdateStatus(ctx context.Context, purchaseID string, status string) error
	GetLeaderboard(
		ctx context.Context,
		filters dto.GetMostPurchases,
	) ([]form.MostPurchases, error)
	FindUserPurchasesByTicket(
		ctx context.Context,
		ticketNumber string,
	) (form.SearchResult, error)
}

type service struct {
	repo       repository.PurchaseRepository
	ticketRepo repository.TicketRepository
	emailer    email.Mailer
}

func NewService(db database.DB, emailClient email.Mailer) Service {
	return &service{
		repo:       repository.NewPurchaseRepository(db),
		ticketRepo: repository.NewTicketRepository(db),
		emailer:    emailClient,
	}
}

func (s *service) Create(
	ctx context.Context,
	req *form.CreatePurchaseRequest,
) error {
	compressedScreenshot, err := utils.CompressToJPG(req.PaymentScreenshot)
	if err != nil {
		return err
	}

	purchase := &types.Purchase{
		UserID:            req.UserID,
		Quantity:          req.Quantity,
		MontoBs:           req.MontoBs,
		MontoUSD:          req.MontoUSD,
		PaymentMethod:     req.PaymentMethod,
		TransactionDigits: req.TransactionDigits,
		PaymentScreenshot: compressedScreenshot,
		Status:            types.StatusPending,
		CreatedAt:         time.Now(),
	}

	purchaseID, err := s.repo.Create(ctx, purchase)
	if err != nil {
		return err
	}

	lotteryID, err := s.ticketRepo.GetActiveLotteryID(ctx)
	if err != nil {
		return err
	}

	_, err = s.ticketRepo.AssignTickets(
		ctx,
		lotteryID,
		req.UserID,
		purchaseID,
		req.SelectedNumbers,
		req.Quantity,
	)
	if err != nil {
		return err
	}

	go func(p *types.Purchase) {
		err := s.emailer.SendPurchaseConfirmation(*p)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("New purchase received and email send! ")
		}
	}(purchase)

	return nil
}

func (s *service) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, int, error) {
	return s.repo.GetAll(ctx, filters)
}

func (s *service) UpdateStatus(
	ctx context.Context,
	purchaseID,
	status string,
) error {
	return s.repo.UpdateStatus(ctx, purchaseID, status)
}

func (s *service) GetLeaderboard(
	ctx context.Context,
	filters dto.GetMostPurchases,
) ([]form.MostPurchases, error) {
	return s.repo.GetLeaderboard(ctx, filters)
}

func (s *service) FindUserPurchasesByTicket(
	ctx context.Context,
	ticketNumber string,
) (form.SearchResult, error) {
	lotteryID, err := s.ticketRepo.GetActiveLotteryID(ctx)
	if err != nil {
		return form.SearchResult{}, err
	}

	user, err := s.repo.FindUserPurchasesByTicket(ctx, lotteryID, ticketNumber)
	if err != nil {
		return form.SearchResult{}, err
	}

	return user, nil
}
