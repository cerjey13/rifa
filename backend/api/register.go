package api

import (
	"rifa/backend/api/httpx"
	"rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHttpRoutes(api huma.API, db db.DB, secureCookies bool) {
	httpx.RegisterAuthRoutes(api, secureCookies, db)
	httpx.RegisterPurchaseRoutes(api, db)
	httpx.RegisterTicketsRoutes(api, db)
	httpx.RegisterPriceRoutes(api, db)
}
