package auth

import (
	"context"
	"errors"

	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"
)

type Service interface {
	Register(ctx context.Context, input *form.RegisterRequest) error
	Login(ctx context.Context, input *form.LoginRequest) (types.AuthUser, error)
}

type service struct {
	users  repository.UserRepository
	config config.ServiceOpts
}

func NewAuthService(db database.DB, opts config.ServiceOpts) Service {
	return &service{
		users:  repository.NewUserRepository(db),
		config: opts,
	}
}

func (s *service) Register(
	ctx context.Context,
	input *form.RegisterRequest,
) error {
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := &types.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: string(hashed),
		Role:     types.CustomerRole,
	}

	err = s.users.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(
	ctx context.Context,
	input *form.LoginRequest,
) (types.AuthUser, error) {
	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return types.AuthUser{}, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		return types.AuthUser{}, errors.New("invalid email or password")
	}

	jwt, err := utils.GenerateJWT(user, s.config.JwtOpts)
	if err != nil {
		return types.AuthUser{}, err
	}
	authUser := types.AuthUser{
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        string(user.Role),
		AccessToken: jwt,
	}

	return authUser, nil
}
