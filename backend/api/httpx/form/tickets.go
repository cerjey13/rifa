package form

type tickets []string

type SearchAvailableTickets struct {
	Unavailable tickets `json:"tickets"`
}

type TicketsPercentage struct {
	Percentage float64 `json:"vendidos"`
}

type UserTickets struct {
	Tickets tickets `json:"tickets"`
}
