package purchase

import (
	"context"
	"errors"
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
	db         database.DB
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
		db:         db,
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
		s.logger.Error(
			ctx,
			"Failed to compress payment image",
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

	tx, txErr := s.db.BeginTx(ctx)
	if txErr != nil {
		s.logger.Error(ctx, "Failed to begin transaction", "error", txErr)
		return txErr
	}
	defer func() {
		if txErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error(
					ctx,
					"Failed to rollback transaction",
					"error",
					rbErr,
				)
				txErr = errors.Join(txErr, rbErr)
			}
		}
	}()

	var purchaseID string
	txErr = func() error {
		purchaseID, err := s.repo.CreateTx(ctx, tx, purchase)
		if err != nil {
			s.logger.Error(
				ctx,
				"Failed to create purchase order",
				"user_id",
				purchase.UserID,
				"error",
				err,
			)
			return err
		}

		lotteryID, err := s.ticketRepo.GetActiveLotteryIDTx(ctx, tx)
		if err != nil {
			s.logger.Error(
				ctx,
				"Failed to get active lottery",
				"user_id",
				purchase.UserID,
				"purchase_id",
				purchaseID,
				"error",
				err,
			)
			return err
		}

		_, err = s.ticketRepo.AssignTicketsTx(
			ctx,
			tx,
			lotteryID,
			purchase.UserID,
			purchaseID,
			req.SelectedNumbers,
			purchase.Quantity,
		)
		if err != nil {
			s.logger.Error(
				ctx,
				"Failed to assign tickets",
				"user_id",
				purchase.UserID,
				"purchase_id",
				purchaseID,
				"error",
				err,
			)
			return err
		}

		return nil
	}()
	if txErr != nil {
		return txErr
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		s.logger.Error(
			ctx,
			"Failed to commit transaction",
			"user_id",
			purchase.UserID,
			"error",
			commitErr,
		)
		return commitErr
	}

	go func(p *types.Purchase, pID string) {
		ct, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err := s.emailer.SendPurchaseConfirmation(ct, *p)
		if err != nil {
			s.logger.Warn(
				ctx,
				"Failed to send the email purchase confirmation",
				"user_id",
				p.UserID,
				"purchase_id",
				pID,
				"error",
				err,
			)
			return
		}

		s.logger.Info(ct, "New purchase received and email send!")
	}(purchase, purchaseID)

	return nil
}

func (s *service) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, int, error) {
	purchases, total, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		s.logger.Error(
			ctx,
			"Failed to get all purchases",
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
			ctx,
			"Failed to update purchase",
			"purchase",
			purchaseID,
			"updated status",
			status,
			"error",
			err,
		)
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
		s.logger.Error(
			ctx,
			"Failed to get purchases leaderboard",
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
		s.logger.Error(ctx, "Failed to get active lottery", "error", err)
		return form.SearchResult{}, err
	}

	user, err := s.repo.FindUserPurchasesByTicket(ctx, lotteryID, ticketNumber)
	if err != nil {
		s.logger.Error(
			ctx,
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
