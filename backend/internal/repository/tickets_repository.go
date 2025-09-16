package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"rifa/backend/internal/types"
	db "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"
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
	GetUnavailableNumbers(
		ctx context.Context,
		lotteryID string,
		numbers []int,
	) ([]int, error)
	GetAvailabilityPercentage(
		ctx context.Context,
		lotteryID string,
	) (float64, error)
	GetUserTickets(
		ctx context.Context,
		userID string,
		lotteryID string,
	) ([]int, error)
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
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("transaction rollback failed: %v", rbErr)
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				log.Printf("transaction commit failed: %v", cmErr)
			}
		}
	}()

	var intNumbers []int
	if len(selectedNumbers) > 0 {
		intNumbers, err = utils.ConvertToIntSlice(selectedNumbers)
		if err != nil {
			return nil, err
		}
	}

	assigned := []types.Ticket{}

	// Assign explicitly chosen numbers
	for _, number := range intNumbers {
		ticket := types.Ticket{}
		res := tx.QueryRow(ctx,
			`UPDATE tickets
			 SET user_id = $1, status = 'sold', purchase_id = $2
			 WHERE lottery_id = $3 AND number = $4 AND status = 'available'
			 RETURNING id, number`,
			userID, purchaseID, lotteryID, number,
		)
		if err := res.Scan(&ticket.ID, &ticket.Number); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf(
					"ticket number %d is no longer available",
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

	// Assign random available tickets for the remainder
	remaining := quantity - len(intNumbers)
	if remaining > 0 {
		rows, err := tx.Query(ctx,
			`SELECT id
			 FROM tickets
			 WHERE lottery_id = $1 AND status = 'available'
			 ORDER BY random()
			 FOR UPDATE SKIP LOCKED
			 LIMIT $2`,
			lotteryID, remaining,
		)
		if err != nil {
			return nil, err
		}
		var ticketIDs []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return nil, err
			}
			ticketIDs = append(ticketIDs, id)
		}
		rows.Close()

		if len(ticketIDs) < remaining {
			return nil, errors.New("not enough tickets available")
		}

		// Assign and return full updated rows
		rows, err = tx.Query(ctx,
			`UPDATE tickets
			 SET user_id = $1, status = 'sold', purchase_id = $2
			 WHERE id = ANY($3)
			 RETURNING id, number`,
			userID, purchaseID, ticketIDs,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var ticket types.Ticket
			if err := rows.Scan(&ticket.ID, &ticket.Number); err != nil {
				return nil, err
			}
			ticket.UserID = &userID
			ticket.Status = "sold"
			ticket.PurchaseID = &purchaseID
			ticket.LotteryID = lotteryID
			assigned = append(assigned, ticket)
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
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

func (r *ticketRepo) GetUnavailableNumbers(
	ctx context.Context,
	lotteryID string,
	numbers []int,
) ([]int, error) {
	if len(numbers) == 0 {
		return nil, nil
	}

	query := `
		SELECT number
		FROM tickets
		WHERE lottery_id = $1
		  AND number = ANY($2)
		  AND status != 'available'
	`

	rows, err := r.db.Query(ctx, query, lotteryID, numbers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	unavailable := []int{}
	for rows.Next() {
		var num int
		if err := rows.Scan(&num); err != nil {
			return nil, err
		}
		unavailable = append(unavailable, num)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return unavailable, nil
}

func (r *ticketRepo) GetAvailabilityPercentage(
	ctx context.Context,
	lotteryID string,
) (float64, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status != 'available')::float / 10000 * 100 AS percent
		FROM tickets
		WHERE lottery_id = $1
	`

	var percent float64
	err := r.db.QueryRow(ctx, query, lotteryID).Scan(&percent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows, return 0% safely
			return 0, nil
		}
		return 0, err
	}

	return percent, nil
}

func (r *ticketRepo) GetUserTickets(
	ctx context.Context,
	userID string,
	lotteryID string,
) ([]int, error) {
	query := `SELECT number
		FROM tickets
		WHERE user_id = $1 AND lottery_id = $2
		ORDER BY number ASC`

	rows, err := r.db.Query(ctx, query, userID, lotteryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := []int{}
	for rows.Next() {
		var num int
		if err := rows.Scan(&num); err != nil {
			return nil, err
		}
		tickets = append(tickets, num)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tickets, nil
}
