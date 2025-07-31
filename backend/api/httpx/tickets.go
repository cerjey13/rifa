package httpx

import (
	"context"
	"log"
	"net/http"
	"strings"

	"rifa/backend/api/httpx/dto"
	"rifa/backend/api/httpx/form"
	mymiddlewares "rifa/backend/api/httpx/middlewares"
	ticket "rifa/backend/internal/core/tickets"
	database "rifa/backend/pkg/db"
	"rifa/backend/pkg/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterTicketsRoutes(api huma.API, db database.DB) {
	srv := ticket.NewService(db)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "ticketsPercentage",
			Method:        http.MethodGet,
			Path:          "/api/tickets/percentage",
			Summary:       "percentage of tickets sold",
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			_ *struct{},
		) (*dto.PercentageOfTicketsSoldOutput, error) {
			percentage, err := srv.GetAvailability(ctx)
			if err != nil {
				log.Println(err)
				return nil, huma.Error500InternalServerError(
					"Failed to get percentage",
				)
			}

			return &dto.PercentageOfTicketsSoldOutput{
				Body: form.TicketsPercentage{Percentage: percentage},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "checkTickets",
			Method:        http.MethodGet,
			Path:          "/api/tickets",
			Summary:       "List if tickets are available",
			Middlewares:   huma.Middlewares{mymiddlewares.RequireSession(api)},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *dto.SearchAvalableTickets,
		) (*dto.SearchAvailableTicketsOutput, error) {
			stringTickets := strings.Split(input.Tickets, ",")
			tickets, err := utils.ConvertToIntSlice(stringTickets)
			if err != nil {
				log.Println(err)
				return nil, huma.Error400BadRequest(
					"Numeros mal formateados",
				)
			}

			unavailableTickets, err := srv.SearchTickets(ctx, tickets)
			if err != nil {
				log.Println(err)
				return nil, huma.Error500InternalServerError(
					"Failed to get purchases",
				)
			}

			strTickets := utils.ConvertToStrSlice(unavailableTickets)

			output := dto.SearchAvailableTicketsOutput{
				Body: &form.SearchAvailableTickets{
					Unavailable: strTickets,
				},
			}
			return &output, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "list user purchases",
			Method:        http.MethodGet,
			Path:          "/api/tickets/users",
			Summary:       "Get purchases for a user",
			Middlewares:   huma.Middlewares{mymiddlewares.RequireSession(api)},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			_ *struct{},
		) (*dto.UserTicketsOutput, error) {
			claims, ok := ctx.Value("claims").(jwt.MapClaims)
			if !ok {
				return nil, huma.Error401Unauthorized("No session claims")
			}

			tickets, err := srv.GetUserTickets(ctx, claims["id"].(string))
			if err != nil {
				return nil, huma.Error500InternalServerError(
					"failed to retrieve user tickets",
				)
			}

			stringTickets := utils.ConvertToStrSlice(tickets)

			return &dto.UserTicketsOutput{
				Body: form.UserTickets{Tickets: stringTickets},
			}, nil
		},
	)
}
