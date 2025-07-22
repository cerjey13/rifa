package repository

import (
	"context"
	"fmt"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type PurchaseRepository interface {
	Create(ctx context.Context, p *types.Purchase) error
	GetAll(
		ctx context.Context,
		filters dto.GetAllPurchases,
	) ([]form.Purchases, error)
}

type purchaseRepo struct{ db database.DB }

func NewPurchaseRepository(db database.DB) PurchaseRepository {
	return &purchaseRepo{db: db}
}

func (r *purchaseRepo) Create(ctx context.Context, p *types.Purchase) error {
	err := r.db.ExecContext(
		ctx,
		`INSERT INTO purchases
		(user_id, quantity, monto_bs, monto_usd, payment_method, 
		transaction_digits, payment_screenshot, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		p.UserID,
		p.Quantity,
		p.MontoBs,
		p.MontoUSD,
		p.PaymentMethod,
		p.TransactionDigits,
		p.PaymentScreenshot,
		p.Status,
		p.CreatedAt,
	)
	return err
}

func (r *purchaseRepo) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, error) {
	var args []interface{}
	query := `SELECT u.name, u.email, u.phone,
    p.quantity, p.monto_bs, p.monto_usd, p.payment_method,
    p.transaction_digits, p.payment_screenshot, p.status, p.created_at
	FROM purchases p JOIN users u ON p.user_id = u.id `

	argIdx := 1
	if filters.PurchaseStatus != "" {
		query += "WHERE P.status = $" + fmt.Sprint(argIdx) + " "
		args = append(args, filters.PurchaseStatus)
		argIdx++
	}

	// Pagination logic
	perPage := filters.ItemCount
	if perPage <= 0 {
		perPage = 10
	}
	page := filters.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage
	query += "ORDER BY p.created_at DESC "
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []form.Purchases
	for rows.Next() {
		var p form.Purchases
		err := rows.Scan(
			&p.User.Name,
			&p.User.Email,
			&p.User.Phone,
			&p.Quantity,
			&p.MontoBs,
			&p.MontoUSD,
			&p.PaymentMethod,
			&p.TransactionDigits,
			&p.PaymentScreenshot,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		purchases = append(purchases, p)
	}

	return purchases, nil
}
