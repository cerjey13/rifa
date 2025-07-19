package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/core/auth"
	database "rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAuthRoutes(api huma.API, db database.DB) {
	svc := auth.NewAuthService(db)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "register",
			Method:        http.MethodPost,
			Path:          "/api/register",
			Summary:       "register a user",
			DefaultStatus: http.StatusCreated,
		},
		func(
			ctx context.Context,
			input *dto.RegisterInput,
		) (*dto.RegisterOutput, error) {
			fmt.Println(input.Body, input)
			err := svc.Register(ctx, &input.Body)
			if err != nil {
				return nil, huma.Error400BadRequest(
					fmt.Sprintf("Registration failed: %v", err),
				)
			}

			return &dto.RegisterOutput{
				Body: form.BaseResponse{Message: "Registration successful"},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "login",
			Method:        http.MethodPost,
			Path:          "/api/login",
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *dto.LoginInput,
		) (*dto.LoginOutput, error) {
			user, err := svc.Login(ctx, &input.Body)
			if err != nil {
				return nil, huma.Error400BadRequest("Invalid credentials")
			}

			return &dto.LoginOutput{
				Body: form.LoginResponse{
					Name:  user.Name,
					Email: user.Email,
					Phone: user.Phone,
				},
				SetCookie: http.Cookie{
					Name:     "session",
					Value:    user.AccessToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
					Expires:  time.Now().Add(7 * 24 * time.Hour),
				},
			}, nil
		},
	)
}
