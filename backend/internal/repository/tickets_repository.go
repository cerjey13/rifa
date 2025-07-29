package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rifa/backend/internal/types"
	db "rifa/backend/pkg/db"
)

type TicketRepository interface {
	AssignTickets(
		ctx context.Context,
		lotteryID,
		userID,
		purchaseID string,
		selectedNumbers []string,
		quantity int,
	) ([]types.Ticket, error)
	GetActiveLotteryID(ctx context.Context) (string, error)
}

type ticketRepo struct {
	db db.DB
}

func NewTicketRepository(db db.DB) TicketRepository {
	return &ticketRepo{db: db}
}

// AssignTickets Assign the selected numbers (if provided and available), and
// assign randoms for the rest
func (r *ticketRepo) AssignTickets(
	ctx context.Context,
	lotteryID,
	userID,
	purchaseID string,
	selectedNumbers []string,
	quantity int,
) ([]types.Ticket, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	assigned := []types.Ticket{}

	// Assign explicitly chosen numbers (safe already)
	for _, number := range selectedNumbers {
		ticket := types.Ticket{}
		res := tx.QueryRow(ctx,
			`UPDATE tickets
			SET user_id = $1, status = 'sold', purchase_id = $2
			WHERE lottery_id = $3 AND number = $4 AND status = 'available'
			RETURNING id, number`,
			userID, purchaseID, lotteryID, number)
		if err := res.Scan(&ticket.ID, &ticket.Number); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf(
					"ticket number %s is no longer available",
					number,
				)
			}
			return nil, err
		}
		ticket.UserID = &userID
		ticket.Status = "sold"
		ticket.PurchaseID = &purchaseID
		ticket.LotteryID = lotteryID
		assigned = append(assigned, ticket)
	}

	// Assign random available tickets for the remainder (now safe)
	remaining := quantity - len(selectedNumbers)
	if remaining > 0 {
		// Select and lock available tickets
		rows, err := tx.Query(ctx,
			`SELECT id, number FROM tickets
			WHERE lottery_id = $1 AND status = 'available'
			FOR UPDATE SKIP LOCKED
			LIMIT $2`,
			lotteryID, remaining)
		if err != nil {
			return nil, err
		}
		var availableTickets []types.Ticket
		var ticketIDs []string
		for rows.Next() {
			var ticket types.Ticket
			if err := rows.Scan(&ticket.ID, &ticket.Number); err != nil {
				return nil, err
			}
			ticketIDs = append(ticketIDs, ticket.ID)
			ticket.UserID = &userID
			ticket.Status = "sold"
			ticket.PurchaseID = &purchaseID
			ticket.LotteryID = lotteryID
			availableTickets = append(availableTickets, ticket)
		}
		rows.Close()
		if len(availableTickets) < remaining {
			return nil, errors.New("not enough tickets available")
		}

		// Batch update the selected tickets
		err = tx.ExecContext(ctx,
			`UPDATE tickets
			SET user_id = $1, status = 'sold', purchase_id = $2
			WHERE id = ANY($3)`,
			userID, purchaseID, ticketIDs)
		if err != nil {
			return nil, err
		}

		assigned = append(assigned, availableTickets...)
	}

	return assigned, nil
}

func (r *ticketRepo) GetActiveLotteryID(ctx context.Context) (string, error) {
	var lotteryID string
	err := r.db.QueryRow(
		ctx,
		`SELECT id FROM lotteries WHERE active = TRUE LIMIT 1`,
	).Scan(&lotteryID)
	return lotteryID, err
}
