package types

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
