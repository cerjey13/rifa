package form

import "github.com/danielgtaylor/huma/v2"

type PurchaseInput struct {
	Form Form
}

type Form struct {
	Quantity          int           `form:"quantity" validate:"min=1"`
	MontoBs           string        `form:"montoBs"`
	MontoUSD          string        `form:"montoUSD"`
	PaymentMethod     string        `form:"paymentMethod" validate:"required"`
	TransactionDigits string        `form:"transactionDigits" validate:"len>=6"`
	PaymentScreenshot huma.FormFile `form:"paymentScreenshot" validate:"required"`
}

type CreatePurchaseRequest struct {
	UserID            string
	Quantity          int
	MontoBs           float64
	MontoUSD          float64
	PaymentMethod     string
	TransactionDigits string
	PaymentScreenshot []byte
}

type Purchases struct {
	UserID            string
	Quantity          int
	MontoBs           float64
	MontoUSD          float64
	PaymentMethod     string
	TransactionDigits string
	PaymentScreenshot []byte
}
