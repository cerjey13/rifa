package purchase

import (
	"context"
	"time"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/core/email"
	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
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
	logger     logx.Logger
	emailer    email.Mailer
}

func NewService(
	db database.DB,
	logger logx.Logger,
	emailClient email.Mailer,
) Service {
	return &service{
		repo:       repository.NewPurchaseRepository(db),
		ticketRepo: repository.NewTicketRepository(db),
		logger:     logger,
		emailer:    emailClient,
	}
}

func (s *service) Create(
	ctx context.Context,
	req *form.CreatePurchaseRequest,
) error {
	compressedScreenshot, err := utils.CompressToJPG(req.PaymentScreenshot)
	if err != nil {
		s.logger.Error("Failed to compress payment image",
			"user_id",
			req.UserID,
			"error",
			err,
		)
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
		s.logger.Error("Failed to create the purchase order",
			"user_id",
			req.UserID,
			"error",
			err,
		)
		return err
	}

	lotteryID, err := s.ticketRepo.GetActiveLotteryID(ctx)
	if err != nil {
		s.logger.Error("Failed to get active lottery",
			"user_id",
			req.UserID,
			"purchase_id",
			purchaseID,
			"error",
			err,
		)
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
		s.logger.Error(
			"Failed to create purchase tickets",
			"user_id",
			req.UserID,
			"purchase_id",
			purchaseID,
			"error",
			err,
		)
		return err
	}

	go func(ctx context.Context, p *types.Purchase) {
		c, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := s.emailer.SendPurchaseConfirmation(c, *p)
		if err != nil {
			s.logger.Warn(
				"Failed to send the email purchase confirmation",
				"user_id",
				req.UserID,
				"purchase_id",
				purchaseID,
				"error",
				err,
			)
			return
		}

		s.logger.Info("New purchase received and email send! ")
	}(ctx, purchase)

	return nil
}

func (s *service) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, int, error) {
	purchases, total, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to get all purchases",
			"page",
			filters.Page,
			"error",
			err,
		)
		return nil, 0, err
	}

	return purchases, total, nil
}

func (s *service) UpdateStatus(
	ctx context.Context,
	purchaseID,
	status string,
) error {
	err := s.repo.UpdateStatus(ctx, purchaseID, status)
	if err != nil {
		s.logger.Error(
			"Failed to update purchase",
			"purchase",
			purchaseID,
			"updated status",
			status,
			"error",
			err)
		return err
	}

	return nil
}

func (s *service) GetLeaderboard(
	ctx context.Context,
	filters dto.GetMostPurchases,
) ([]form.MostPurchases, error) {
	leaderboard, err := s.repo.GetLeaderboard(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to get purchases leaderboard",
			"page",
			filters.Page,
			"error",
			err,
		)
		return nil, err
	}

	return leaderboard, nil
}

func (s *service) FindUserPurchasesByTicket(
	ctx context.Context,
	ticketNumber string,
) (form.SearchResult, error) {
	lotteryID, err := s.ticketRepo.GetActiveLotteryID(ctx)
	if err != nil {
		s.logger.Error("Failed to get active lottery", "error", err)
		return form.SearchResult{}, err
	}

	user, err := s.repo.FindUserPurchasesByTicket(ctx, lotteryID, ticketNumber)
	if err != nil {
		s.logger.Error(
			"Failed to find user purchases by ticket number",
			"ticket",
			ticketNumber,
			"error",
			err,
		)
		return form.SearchResult{}, err
	}

	return user, nil
}
