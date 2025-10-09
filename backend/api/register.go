package api

import (
	"rifa/backend/api/httpx"
	"rifa/backend/pkg/config"
	"rifa/backend/pkg/db"
	"rifa/backend/pkg/logx"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterHttpRoutes(
	api huma.API,
	db db.DB,
	logger logx.Logger,
	serviceOpts config.ServiceOpts,
) {
	httpx.RegisterAuthRoutes(api, db, logger, serviceOpts)
	httpx.RegisterPurchaseRoutes(api, db, logger, serviceOpts)
	httpx.RegisterTicketsRoutes(api, db, logger, serviceOpts)
	httpx.RegisterPriceRoutes(api, db, logger, serviceOpts)
}
