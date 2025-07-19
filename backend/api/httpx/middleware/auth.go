package middleware

import (
	"net/http"
	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
)

func RequireSession(ctx huma.Context, next func(huma.Context)) {
	cookie, err := huma.ReadCookie(ctx, "session")
	if err != nil || cookie == nil || cookie.Value == "" {
		_ = huma.WriteErr(
			nil,
			ctx,
			http.StatusUnauthorized,
			"Unauthenticated",
			nil,
		)
		return
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		_ = huma.WriteErr(
			nil,
			ctx,
			http.StatusUnauthorized,
			"Invalid or expired session",
			nil,
		)
		return
	}

	ctx = huma.WithValue(ctx, "claims", claims)
	next(ctx)
}
