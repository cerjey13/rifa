package dto

import "rifa/backend/api/httpx/form"

type PriceOutput struct {
	Body form.LotteryPrices
}

type PriceUpdateInput struct {
	Body form.LotteryPrices
}
