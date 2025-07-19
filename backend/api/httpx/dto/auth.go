package dto

import (
	"net/http"

	"rifa/backend/api/httpx/form"
)

type RegisterInput struct {
	Body form.RegisterRequest
}

type RegisterOutput struct {
	Body form.BaseResponse
}

type LoginInput struct {
	Body form.LoginRequest
}

type LoginOutput struct {
	Body      form.BaseResponse
	SetCookie http.Cookie `header:"Set-Cookie"`
}
