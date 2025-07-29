package form

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type PurchaseInput struct {
	Form Form
}

type Form struct {
	Quantity          int           `form:"quantity" validate:"min=2"`
	MontoBs           string        `form:"montoBs"`
	MontoUSD          string        `form:"montoUSD"`
	PaymentMethod     string        `form:"paymentMethod" validate:"required"`
	TransactionDigits string        `form:"transactionDigits" validate:"len>=6"`
	SelectedNumbers   string        `form:"selectedNumbers"`
	PaymentScreenshot huma.FormFile `form:"paymentScreenshot" validate:"required"`
}

type CreatePurchaseRequest struct {
	UserID            string
	Quantity          int
	MontoBs           float64
	MontoUSD          float64
	PaymentMethod     string
	TransactionDigits string
	SelectedNumbers   []string
	PaymentScreenshot []byte
}

type Purchases struct {
	User              User      `json:"user"`
	ID                string    `json:"id"`
	Quantity          int       `json:"quantity"`
	MontoBs           float64   `json:"montoBs"`
	MontoUSD          float64   `json:"montoUsd"`
	PaymentMethod     string    `json:"paymentMethod"`
	TransactionDigits string    `json:"transactionDigits"`
	Status            string    `json:"status"`
	PaymentScreenshot []byte    `json:"paymentScreenshot"`
	CreatedAt         time.Time `json:"date"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type MostPurchases struct {
	User     User `json:"user"`
	Quantity int  `json:"quantity"`
}
