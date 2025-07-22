package purchase

import (
	"context"
	"fmt"
	"time"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type Service interface {
	Create(ctx context.Context, req *form.CreatePurchaseRequest) error
	GetAll(
		ctx context.Context,
		filters dto.GetAllPurchases,
	) ([]form.Purchases, error)
}

type service struct {
	repo repository.PurchaseRepository
}

func NewService(db database.DB) Service {
	return &service{repo: repository.NewPurchaseRepository(db)}
}

func (s *service) Create(
	ctx context.Context,
	req *form.CreatePurchaseRequest,
) error {
	purchase := &types.Purchase{
		UserID:            req.UserID,
		Quantity:          req.Quantity,
		MontoBs:           req.MontoBs,
		MontoUSD:          req.MontoUSD,
		PaymentMethod:     req.PaymentMethod,
		TransactionDigits: req.TransactionDigits,
		PaymentScreenshot: req.PaymentScreenshot,
		Status:            types.StatusPending,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.Create(ctx, purchase); err != nil {
		return err
	}

	// TODO: Send confirmation email to admin with purchase data (use Mailgun or similar).
	// For now, mock/send to logs:
	go func(p *types.Purchase) {
		// Here you would integrate with Mailgun/sendgrid/SMTP etc.
		// e.g., sendPurchaseEmail(p)
		fmt.Printf("Mock email: New purchase received! %+v\n", p)
	}(purchase)

	return nil
}

func (s *service) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, error) {
	return s.repo.GetAll(ctx, filters)
}
