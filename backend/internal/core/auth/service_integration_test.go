//go:build integration
// +build integration

package auth

import (
	"context"
	"testing"
	"time"

	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/types"
	"rifa/backend/pkg/config"
	"rifa/backend/pkg/logx"
	"rifa/backend/testutils"
)

func TestAuthService_Integration(t *testing.T) {
	pg := testutils.StartPostgres(t)
	defer testutils.StopPostgres(t, pg)

	logger := logx.NewLogger("development")
	cfg := config.ServiceOpts{
		JwtOpts: config.JwtOpts{
			JwtSecret:    "testsecret",
			JwtExpiresAt: 3600,
		},
	}

	svc := NewAuthService(pg.Database, logger, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		name        string
		registerReq *form.RegisterRequest
		loginReq    *form.LoginRequest
		expectErr   bool
		validate    func(t *testing.T, got types.AuthUser, err error)
	}{
		{
			name: "Register new user and login successfully",
			registerReq: &form.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "123456789",
				Password: "pass1234",
			},
			loginReq: &form.LoginRequest{
				Email:    "john@example.com",
				Password: "pass1234",
			},
			expectErr: false,
			validate: func(t *testing.T, got types.AuthUser, err error) {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if got.Email != "john@example.com" {
					t.Errorf("expected email %q, got %q", "john@example.com", got.Email)
				}
				if got.AccessToken == "" {
					t.Error("expected non-empty JWT token")
				}
			},
		},
		{
			name: "Fail to login with wrong password",
			registerReq: &form.RegisterRequest{
				Name:     "Jane Doe",
				Email:    "jane@example.com",
				Phone:    "987654321",
				Password: "strongpass",
			},
			loginReq: &form.LoginRequest{
				Email:    "jane@example.com",
				Password: "wrongpass",
			},
			expectErr: true,
			validate: func(t *testing.T, _ types.AuthUser, err error) {
				if err == nil {
					t.Error("expected error for invalid password, got nil")
				}
			},
		},
		{
			name: "Fail duplicate registration",
			registerReq: &form.RegisterRequest{
				Name:     "John Dup",
				Email:    "dup@example.com",
				Phone:    "999999999",
				Password: "dupPass123",
			},
			loginReq: &form.LoginRequest{
				Email:    "dup@example.com",
				Password: "dupPass123",
			},
			expectErr: false,
			validate: func(t *testing.T, got types.AuthUser, err error) {
				if err != nil {
					t.Fatalf("expected first registration to succeed: %v", err)
				}

				// Try registering again with the same email
				dupReq := &form.RegisterRequest{
					Name:     "Dup Again",
					Email:    "dup@example.com",
					Phone:    "888888888",
					Password: "dupPass123",
				}
				errDup := svc.Register(context.Background(), dupReq)
				if errDup == nil {
					t.Error("expected duplicate registration to fail, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.registerReq != nil {
				if err := svc.Register(ctx, tt.registerReq); err != nil {
					if tt.expectErr {
						t.Logf("expected error during register: %v", err)
						return
					}
					t.Fatalf("unexpected register error: %v", err)
				}
			}

			var (
				got types.AuthUser
				err error
			)
			if tt.loginReq != nil {
				got, err = svc.Login(ctx, tt.loginReq)
			}

			tt.validate(t, got, err)
		})
	}
}
