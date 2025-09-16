package types

type Ticket struct {
	ID         string
	LotteryID  string
	Number     string
	UserID     *string
	Status     string
	PurchaseID *string
}
