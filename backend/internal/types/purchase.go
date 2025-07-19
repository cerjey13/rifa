package types

import "time"

type PurchaseStatus string

const (
	StatusPending   PurchaseStatus = "pending"
	StatusVerified  PurchaseStatus = "verified"
	StatusCancelled PurchaseStatus = "cancelled"
)

type Purchase struct {
	UserID            string
	Quantity          int
	MontoBs           float64
	MontoUSD          float64
	PaymentMethod     string
	TransactionDigits string
	PaymentScreenshot []byte
	Status            PurchaseStatus
	CreatedAt         time.Time
}
