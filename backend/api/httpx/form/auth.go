package form

type RegisterRequest struct {
	Name     string `json:"name" required:"true"`
	Email    string `json:"email" required:"true"`
	Phone    string `json:"phone" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginRequest struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type BaseResponse struct {
	Message string `json:"message"`
}
