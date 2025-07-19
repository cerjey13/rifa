package repository

import (
	"context"
	"fmt"

	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type PurchaseRepository interface {
	Create(ctx context.Context, p *types.Purchase) error
	GetAll(ctx context.Context) ([]*types.Purchase, error)
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
	fmt.Println(err)
	return err
}

func (r *purchaseRepo) GetAll(ctx context.Context) ([]*types.Purchase, error) {
	rows, err := r.db.Query(ctx, `
        SELECT user_id, quantity, monto_bs, monto_usd, payment_method, transaction_digits, payment_screenshot, status, created_at
        FROM purchases
        ORDER BY created_at DESC
    `)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var purchases []*types.Purchase
	for rows.Next() {
		var p types.Purchase
		err := rows.Scan(
			&p.UserID,
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
			fmt.Println(err)
			return nil, err
		}
		purchases = append(purchases, &p)
	}
	return purchases, nil
}
