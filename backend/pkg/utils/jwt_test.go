package utils

import (
	"strings"
	"testing"
	"time"

	"rifa/backend/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "unit-test-secret"

// -------- helpers --------

func mustClaimString(t *testing.T, claims jwt.MapClaims, key string) string {
	t.Helper()
	v, ok := claims[key]
	if !ok {
		t.Fatalf("claim %q missing", key)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("claim %q not a string: %#v", key, v)
	}
	return s
}

func mustClaimInt64(t *testing.T, claims jwt.MapClaims, key string) int64 {
	t.Helper()
	v, ok := claims[key]
	if !ok {
		t.Fatalf("claim %q missing", key)
	}
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	default:
		t.Fatalf("claim %q not a number: %#v", key, v)
		return 0
	}
}

// -------- tests --------

func TestHashPassword_And_CheckPassword(t *testing.T) {
	const pwd = "Str0ngP@ss!"
	hash, err := HashPassword(pwd)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Fatalf("empty hash")
	}
	if hash == pwd {
		t.Fatalf("hash should not equal plaintext")
	}
	// bcrypt hashes start with $2 (e.g., $2a$, $2b$, $2y$)
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("unexpected hash prefix: %q", hash[:2])
	}

	if !CheckPassword(pwd, hash) {
		t.Fatalf("CheckPassword should succeed for correct password")
	}
	if CheckPassword("wrongpassword", hash) {
		t.Fatalf("CheckPassword should fail for wrong password")
	}
}

func TestGenerateJWT_And_ValidateJWT_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	u := &types.User{
		ID:    "user-123",
		Name:  "Fran",
		Email: "fran@example.com",
		Role:  types.AdminRole,
	}

	tok, err := GenerateJWT(u)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}
	if tok == "" {
		t.Fatalf("empty token")
	}

	claims, err := ValidateJWT(tok)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}

	if got := mustClaimString(t, claims, "id"); got != u.ID {
		t.Fatalf("id claim = %q, want %q", got, u.ID)
	}
	if got := mustClaimString(t, claims, "name"); got != u.Name {
		t.Fatalf("name claim = %q, want %q", got, u.Name)
	}
	if got := mustClaimString(t, claims, "email"); got != u.Email {
		t.Fatalf("email claim = %q, want %q", got, u.Email)
	}
	if got := mustClaimString(t, claims, "role"); got != string(u.Role) {
		t.Fatalf("role claim = %q, want %q", got, u.Role)
	}

	exp := mustClaimInt64(t, claims, "exp")
	if exp <= time.Now().Unix() {
		t.Fatalf("exp claim not in the future: %d", exp)
	}
}

func TestValidateJWT_Failures(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	now := time.Now()

	// Build tokens for failure scenarios.
	makeSigned := func(secret string, exp time.Time) string {
		claims := jwt.MapClaims{
			"id":    "u1",
			"name":  "n",
			"email": "e@example.com",
			"role":  "user",
			"exp":   exp.Unix(),
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, err := tok.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("sign token: %v", err)
		}
		return s
	}

	validWrongSig := makeSigned("other-secret", now.Add(1*time.Hour))
	expired := makeSigned(testSecret, now.Add(-1*time.Hour))

	tests := []struct {
		name  string
		token string
	}{
		{"garbage_string", "this-is-not-a-jwt"},
		{"wrong_signature", validWrongSig},
		{"expired_token", expired},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateJWT(tt.token); err == nil {
				t.Fatalf("expected validation failure for case %q", tt.name)
			}
		})
	}
}
