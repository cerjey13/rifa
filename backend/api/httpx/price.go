package httpx

import (
	"context"
	"log"
	"net/http"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/internal/core/price"
	database "rifa/backend/pkg/db"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterPriceRoutes(api huma.API, db database.DB) {
	srv := price.NewService(db)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "prices",
			Method:        http.MethodGet,
			Path:          "/api/prices",
			Summary:       "get the prices amounts",
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			_ *struct{},
		) (*dto.PriceOutput, error) {
			prices, err := srv.GetPrices(ctx)
			if err != nil {
				log.Println(err)
				return nil, huma.Error500InternalServerError(
					"Failed to get prices",
				)
			}

			return &dto.PriceOutput{
				Body: form.LotteryPrices{
					BS:  prices.BsAmount,
					USD: prices.UsdAmount,
				},
			}, nil
		},
	)
}
