package api

import (
	"rifa/backend/api/httpx"
	"rifa/backend/pkg/config"
	"rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHttpRoutes(api huma.API, db db.DB, serviceOpts config.ServiceOpts) {
	httpx.RegisterAuthRoutes(api, db, serviceOpts)
	httpx.RegisterPurchaseRoutes(api, db, serviceOpts)
	httpx.RegisterTicketsRoutes(api, db, serviceOpts)
	httpx.RegisterPriceRoutes(api, db, serviceOpts)
}
