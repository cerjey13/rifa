package dto

import "rifa/backend/api/httpx/form"

type SearchAvailableTicketsOutput struct {
	Body *form.SearchAvailableTickets
}

type SearchAvalableTickets struct {
	Tickets string `query:"numbers"`
}

type PercentageOfTicketsSoldOutput struct {
	Body form.TicketsPercentage
}

type UserTicketsOutput struct {
	Body form.UserTickets
}
