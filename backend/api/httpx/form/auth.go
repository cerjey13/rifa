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

type LoginResponse struct {
	Name  string `json:"name" required:"true"`
	Email string `json:"email" required:"true"`
	Phone string `json:"phone" required:"true"`
	Role  string `json:"role" required:"true"`
}

type BaseResponse struct {
	Message string `json:"message"`
}

type MeResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
