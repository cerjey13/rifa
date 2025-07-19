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
	Body      form.LoginResponse
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type MeOutput struct {
	Body form.MeResponse
}

type LogoutOutput struct {
	Body        form.BaseResponse
	ClearCookie http.Cookie `header:"Set-Cookie"`
}
