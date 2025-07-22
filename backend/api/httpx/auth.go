package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	mymiddlewares "rifa/backend/api/httpx/middlewares"
	"rifa/backend/internal/core/auth"
	database "rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
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
					Role:  user.Role,
				},
				SetCookie: http.Cookie{
					Name:     "session",
					Value:    user.AccessToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
					Expires:  time.Now().Add(7 * 24 * time.Hour),
				},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "me",
			Method:        http.MethodGet,
			Path:          "/api/me",
			Summary:       "Get current user",
			Middlewares:   huma.Middlewares{mymiddlewares.RequireSession(api)},
			DefaultStatus: http.StatusOK,
		},
		func(ctx context.Context, _ *struct{}) (*dto.MeOutput, error) {
			claims, ok := ctx.Value("claims").(jwt.MapClaims)
			if !ok {
				return nil, huma.Error401Unauthorized("No session claims")
			}
			output := &dto.MeOutput{}
			output.Body = form.MeResponse{
				Name:  claims["name"].(string),
				Email: claims["email"].(string),
				Role:  claims["role"].(string),
			}

			return output, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "logout",
			Method:      http.MethodPost,
			Path:        "/api/logout",
			Summary:     "Logout user (expire session cookie)",
		},
		func(ctx context.Context, _ *struct{}) (*dto.LogoutOutput, error) {
			return &dto.LogoutOutput{
				ClearCookie: http.Cookie{
					Name:     "session",
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
					Expires:  time.Now().Add(-1 * time.Hour),
					MaxAge:   -1,
				},
				Body: form.BaseResponse{
					Message: "Logged out",
				},
			}, nil
		},
	)
}
