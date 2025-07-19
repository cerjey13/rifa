package repository

import (
	"context"
	"errors"

	"rifa/backend/internal/types"
	database "rifa/backend/pkg/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetByEmail(ctx context.Context, email string) (*types.User, error)
}

type userRepo struct{ db database.DB }

func NewUserRepository(db database.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *types.User) error {
	query := `
		INSERT INTO users (name, email, phone, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING null
	`
	return r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Role,
	).Scan(&user.ID)
}

func (r *userRepo) GetByEmail(
	ctx context.Context,
	email string,
) (*types.User, error) {
	query := `
	SELECT id, name, email, phone, password, role FROM users WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email)

	var user types.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
