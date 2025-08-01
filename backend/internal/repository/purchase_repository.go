package repository

import (
	"context"
	"fmt"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"
)

type PurchaseRepository interface {
	Create(ctx context.Context, p *types.Purchase) (string, error)
	GetAll(
		ctx context.Context,
		filters dto.GetAllPurchases,
	) ([]form.Purchases, int, error)
	UpdateStatus(ctx context.Context, purchaseID, status string) error
	GetLeaderboard(
		ctx context.Context,
		filters dto.GetMostPurchases,
	) ([]form.MostPurchases, error)
}

type purchaseRepo struct{ db database.DB }

func NewPurchaseRepository(db database.DB) PurchaseRepository {
	return &purchaseRepo{db: db}
}

func (r *purchaseRepo) Create(
	ctx context.Context,
	p *types.Purchase,
) (string, error) {
	var id string
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO purchases
		(user_id, quantity, monto_bs, monto_usd, payment_method, 
		transaction_digits, payment_screenshot, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id`,
		p.UserID,
		p.Quantity,
		p.MontoBs,
		p.MontoUSD,
		p.PaymentMethod,
		p.TransactionDigits,
		p.PaymentScreenshot,
		p.Status,
		p.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *purchaseRepo) GetAll(
	ctx context.Context,
	filters dto.GetAllPurchases,
) ([]form.Purchases, int, error) {
	var args []interface{}
	query := `SELECT u.id, u.name, u.email, u.phone,
    p.id, p.quantity, p.monto_bs, p.monto_usd, p.payment_method,
    p.transaction_digits, p.payment_screenshot, p.status, p.created_at,
	COALESCE(
  	ARRAY_AGG(t.number ORDER BY t.number) FILTER (WHERE t.number IS NOT NULL),
  	'{}') AS numbers,
	COUNT(*) OVER() as total_count
	FROM purchases p JOIN users u ON p.user_id = u.id
	LEFT JOIN tickets t ON t.purchase_id = p.id `

	argIdx := 1
	if filters.PurchaseStatus != "" {
		query += "WHERE p.status = $" + fmt.Sprint(argIdx) + " "
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
	query += `
	GROUP BY 
		u.id, u.name, u.email, u.phone,
		p.id, p.quantity, p.monto_bs, p.monto_usd, p.payment_method,
		p.transaction_digits, p.payment_screenshot, p.status, p.created_at
	ORDER BY p.created_at DESC
	`
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var purchases []form.Purchases
	var total int
	for rows.Next() {
		var p form.Purchases
		var numbers []int
		var rowTotal int
		err := rows.Scan(
			&p.User.ID,
			&p.User.Name,
			&p.User.Email,
			&p.User.Phone,
			&p.ID,
			&p.Quantity,
			&p.MontoBs,
			&p.MontoUSD,
			&p.PaymentMethod,
			&p.TransactionDigits,
			&p.PaymentScreenshot,
			&p.Status,
			&p.CreatedAt,
			&numbers,
			&rowTotal,
		)
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			total = rowTotal
		}
		p.Tickets = utils.ConvertToStrSlice(numbers)
		purchases = append(purchases, p)
	}
	return purchases, total, nil
}

func (r *purchaseRepo) UpdateStatus(
	ctx context.Context,
	purchaseID,
	status string,
) error {
	err := r.db.ExecContext(ctx,
		`UPDATE purchases SET status = $1 WHERE id = $2`,
		status, purchaseID,
	)
	return err
}

func (r *purchaseRepo) GetLeaderboard(
	ctx context.Context,
	filters dto.GetMostPurchases,
) ([]form.MostPurchases, error) {
	perPage := filters.ItemCount
	if perPage <= 0 {
		perPage = 10
	}
	page := filters.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage

	query := `
		SELECT id, name, email, phone, quantity
		FROM (
			SELECT u.id, u.name, u.email, u.phone,
				SUM(p.quantity) as quantity
			FROM purchases p
			JOIN users u ON p.user_id = u.id
			WHERE p.status = 'verified'
			GROUP BY u.id, u.name, u.email, u.phone
		) AS leaderboard
		ORDER BY quantity DESC
		LIMIT $1 OFFSET $2
	`
	args := []interface{}{perPage, offset}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []form.MostPurchases
	for rows.Next() {
		var entry form.MostPurchases
		err := rows.Scan(
			&entry.User.ID,
			&entry.User.Name,
			&entry.User.Email,
			&entry.User.Phone,
			&entry.Quantity,
		)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}
