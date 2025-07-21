package dto

import (
	"rifa/backend/api/httpx/form"

	"github.com/danielgtaylor/huma/v2"
)

type PurchaseInput struct {
	RawBody huma.MultipartFormFiles[struct {
		ScreenShot huma.FormFile `form:"paymentScreenshot" contentType:"image/*" required:"true"`
		form.Form
	}]
}

type PurchaseOutput struct {
	Body form.BaseResponse
}

type PurchasesOutput struct {
	Body []form.Purchases
}
