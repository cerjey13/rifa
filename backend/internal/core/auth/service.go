package auth

import (
	"context"
	"errors"

	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/repository"
	"rifa/backend/internal/types"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"
	"rifa/backend/pkg/utils"
)

type Service interface {
	Register(ctx context.Context, input *form.RegisterRequest) error
	Login(ctx context.Context, input *form.LoginRequest) (types.AuthUser, error)
}

type service struct {
	users  repository.UserRepository
	logger logx.Logger
	config config.ServiceOpts
}

func NewAuthService(
	db database.DB,
	logger logx.Logger,
	opts config.ServiceOpts,
) Service {
	return &service{
		users:  repository.NewUserRepository(db),
		logger: logger,
		config: opts,
	}
}

func (s *service) Register(
	ctx context.Context,
	input *form.RegisterRequest,
) error {
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		s.logger.Error("Failed to hash the password provided", "error", err)
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
		s.logger.Error(
			"Failed to create user",
			"email",
			user.Email,
			"error",
			err,
		)
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
		s.logger.Error(
			"Failed to query user for login",
			"email",
			input.Email,
			"error",
			err,
		)
		return types.AuthUser{}, err
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		s.logger.Warn("Invalid login attempt", "email", input.Email)
		return types.AuthUser{}, errors.New("invalid email or password")
	}

	jwt, err := utils.GenerateJWT(user, s.config.JwtOpts)
	if err != nil {
		s.logger.Error("Failed to generate jwt", "error", err)
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
