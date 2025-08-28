package httpx

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	mymiddlewares "rifa/backend/api/httpx/middlewares"
	"rifa/backend/internal/core/email"
	"rifa/backend/internal/core/purchase"
	"rifa/backend/pkg/config"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

const maxFileBytes int64 = 1 * 1024 * 1024

func RegisterPurchaseRoutes(api huma.API, db database.DB) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config for purchase routes")
	}

	emailer := email.NewMailerooClient(
		cfg.MailerooApiKey,
		cfg.EmailSender,
		cfg.EmailReciever,
		cfg.EmailURL,
	)
	srv := purchase.NewService(db, emailer)

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
			claims, ok := ctx.Value("claims").(jwt.MapClaims)
			if !ok {
				log.Println("No session claims")
				return nil, huma.Error401Unauthorized("No session claims")
			}

			formData := input.RawBody.Data()
			if formData.ScreenShot.Size > maxFileBytes {
				log.Println("purchase image too large ", formData.ScreenShot.Size)
				return nil, huma.NewError(
					http.StatusRequestEntityTooLarge,
					http.StatusText(http.StatusRequestEntityTooLarge),
					errors.New("image exceeds 3MB limit"),
				)
			}

			// Defensive read guard in case Size is missing/wrong
			r := io.LimitReader(formData.ScreenShot, maxFileBytes+1)
			screenshot, err := io.ReadAll(r)
			if err != nil {
				log.Println(err)
				return nil, huma.Error400BadRequest(
					"Could not read uploaded file",
				)
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
				SelectedNumbers: func() []string {
					if formData.SelectedNumbers != "" {
						return strings.Split(formData.SelectedNumbers, ",")

					}
					return []string{}
				}(),
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
		OperationID:   "leaderboard",
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

	huma.Register(
		api,
		huma.Operation{
			OperationID: "updatePurchase",
			Method:      http.MethodPatch,
			Path:        "/api/purchases",
			Summary:     "Update a purchase status (admin only)",
			Middlewares: huma.Middlewares{
				mymiddlewares.RequireAdminSession(api),
			},
			DefaultStatus: http.StatusNoContent,
		},
		func(ctx context.Context, input *dto.UpdatePurchase) (*struct{}, error) {
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

	huma.Register(
		api,
		huma.Operation{
			OperationID: "purchaseByNumber",
			Method:      http.MethodGet,
			Path:        "/api/purchases/search",
			Summary:     "Search for a user data by number bought (admin only)",
			Middlewares: huma.Middlewares{
				mymiddlewares.RequireAdminSession(api),
			},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *dto.SearchPurchase,
		) (*dto.SearchPurchaseOutput, error) {
			searched, err := srv.FindUserPurchasesByTicket(ctx, input.Number)
			if err != nil {
				log.Println(err)
				return nil, huma.Error500InternalServerError(
					"Failed to retrieve number data",
				)
			}

			return &dto.SearchPurchaseOutput{Body: searched}, nil
		},
	)
}
