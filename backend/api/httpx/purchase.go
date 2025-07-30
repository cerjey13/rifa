package httpx

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	mymiddlewares "rifa/backend/api/httpx/middlewares"
	"rifa/backend/internal/core/purchase"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterPurchaseRoutes(api huma.API, db database.DB) {
	srv := purchase.NewService(db)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "purchase",
			Method:        http.MethodPost,
			Path:          "/api/purchases",
			Summary:       "Submit a purchase",
			Middlewares:   huma.Middlewares{mymiddlewares.RequireSession(api)},
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
				SelectedNumbers:   strings.Split(formData.SelectedNumbers, ","),
				PaymentScreenshot: screenshot,
			}
			if err := srv.Create(ctx, purchase); err != nil {
				log.Println(err)
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
		Middlewares:   huma.Middlewares{mymiddlewares.RequireAdminSession(api)},
		DefaultStatus: http.StatusOK,
	}, func(
		ctx context.Context,
		input *dto.GetAllPurchases,
	) (*dto.PurchasesOutput, error) {
		purchases, total, err := srv.GetAll(ctx, *input)
		if err != nil {
			log.Println(err)
			return nil, huma.Error500InternalServerError(
				"Failed to get purchases",
			)
		}

		output := dto.PurchasesOutput{Body: purchases, Total: total}
		return &output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "listPurchases",
		Method:        http.MethodGet,
		Path:          "/api/purchases/leaderboard",
		Summary:       "List purchases by user with the most buyed",
		Middlewares:   huma.Middlewares{mymiddlewares.RequireAdminSession(api)},
		DefaultStatus: http.StatusOK,
	}, func(
		ctx context.Context,
		input *dto.GetMostPurchases,
	) (*dto.MostPurchasesOutput, error) {
		leaderboard, err := srv.GetLeaderboard(ctx, *input)
		if err != nil {
			log.Println(err)
			return nil, huma.Error500InternalServerError(
				"Failed to get purchases",
			)
		}

		output := dto.MostPurchasesOutput{}
		output.Body = leaderboard
		return &output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "list user purchases",
		Method:        http.MethodGet,
		Path:          "/api/purchases/tickets",
		Summary:       "Get purchases for a user",
		Middlewares:   huma.Middlewares{mymiddlewares.RequireAdminSession(api)},
		DefaultStatus: http.StatusOK,
	}, func(
		ctx context.Context,
		_ *struct{},
	) (*struct {
		Body form.MostPurchases
	}, error) {
		claims, ok := ctx.Value("claims").(jwt.MapClaims)
		if !ok {
			return nil, huma.Error401Unauthorized("No session claims")
		}
		//TODO: implement the retrieval of tickets
		output := &struct{ Body form.MostPurchases }{}
		output.Body.Quantity = 100
		output.Body.User.ID = claims["id"].(string)
		return output, nil
	})

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "updatePurchase",
			Method:        http.MethodPatch,
			Path:          "/api/purchases",
			Summary:       "Update a purchase status (admin only)",
			Middlewares:   huma.Middlewares{mymiddlewares.RequireSession(api)},
			DefaultStatus: http.StatusNoContent,
		},
		func(
			ctx context.Context,
			input *dto.UpdatePurchase,
		) (*struct{}, error) {
			err := srv.UpdateStatus(ctx, input.ID, input.Body.Status)
			if err != nil {
				log.Println(err)
				return nil, huma.Error500InternalServerError(
					"Failed to update purchase",
				)
			}

			return nil, nil
		},
	)
}
