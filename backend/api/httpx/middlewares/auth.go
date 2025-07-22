package mymiddlewares

import (
	"log"
	"net/http"

	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
)

func RequireSession(
	api huma.API,
) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		cookie, err := huma.ReadCookie(ctx, "session")
		if err != nil || cookie == nil || cookie.Value == "" {
			log.Println(err)
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusUnauthorized,
				"Unauthenticated",
				err,
			)
			return
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			log.Println(err)
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusUnauthorized,
				"Invalid or expired session",
				err,
			)
			return
		}
		ctx = huma.WithValue(ctx, "claims", claims)
		next(ctx)
	}
}
