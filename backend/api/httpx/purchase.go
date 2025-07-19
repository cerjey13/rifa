package httpx

import (
	"context"
	"io"
	"net/http"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	"rifa/backend/api/httpx/middleware"
	"rifa/backend/internal/core/purchase"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterPurchaseRoutes(api huma.API, db database.DB) {
	svr := purchase.NewService(db)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "purchase",
			Method:        http.MethodPost,
			Path:          "/api/purchase",
			Summary:       "Submit a purchase",
			Middlewares:   huma.Middlewares{middleware.RequireSession(api)},
			DefaultStatus: http.StatusCreated,
		},
		func(
			ctx context.Context,
			input *dto.PurchaseInput,
		) (*dto.PurchaseOutput, error) {
			formData := input.RawBody.Data()
			screenshot, err := io.ReadAll(formData.ScreenShot)
			if err != nil {
				return nil, huma.Error500InternalServerError(
					"Could not read uploaded file",
				)
			}
			claims, ok := ctx.Value("claims").(jwt.MapClaims)
			if !ok {
				return nil, huma.Error401Unauthorized("No session claims")
			}
			// Save in DB
			purchase := &form.CreatePurchaseRequest{
				UserID:   claims["id"].(string),
				Quantity: formData.Quantity,
				MontoBs: utils.ParseFloatOrZero(
					formData.MontoBs,
				),
				MontoUSD: utils.ParseFloatOrZero(
					formData.MontoUSD,
				),
				PaymentMethod:     formData.PaymentMethod,
				TransactionDigits: formData.TransactionDigits,
				PaymentScreenshot: screenshot,
			}
			if err := svr.Create(ctx, purchase); err != nil {
				return nil, huma.Error500InternalServerError(
					"Failed to save purchase",
				)
			}

			return &dto.PurchaseOutput{
				Body: form.BaseResponse{Message: "success"},
			}, nil
		},
	)

	huma.Register(api, huma.Operation{
		OperationID:   "listPurchases",
		Method:        http.MethodGet,
		Path:          "/api/purchases",
		Summary:       "List all purchases (admin only)",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, _ *struct{}) (*dto.PurchasesOutput, error) {
		list, err := svr.GetAll(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to get purchases")
		}
		purchases := make([]form.Purchases, 0, len(list))
		for _, p := range list {
			purchases = append(purchases, form.Purchases{
				UserID:            p.UserID,
				Quantity:          p.Quantity,
				MontoBs:           p.MontoBs,
				MontoUSD:          p.MontoUSD,
				PaymentMethod:     p.PaymentMethod,
				TransactionDigits: p.TransactionDigits,
				PaymentScreenshot: p.PaymentScreenshot,
			})
		}

		output := dto.PurchasesOutput{}
		output.Body = purchases
		return &output, nil
	})

}
