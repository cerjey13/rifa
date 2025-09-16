package config

import (
	"os"
	"sync"
	"testing"
)

// reset clears the package-level singletons so each test runs fresh.
func reset() {
	once = sync.Once{}
	cfg = &Config{}
}

func TestNewConfig_Defaults(t *testing.T) {
	reset()
	t.Cleanup(reset)

	// Ensure env is clean
	err := os.Unsetenv("PORT")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("HOST")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("APP_ENV")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("COOKIE_SECURE")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("JWT_SECRET")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("DATABASE_URL")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("EMAIL_MAILEROO_API_KEY")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("EMAIL_ACCOUNT")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("EMAIL_SENDER_ACCOUNT")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}
	err = os.Unsetenv("EMAIL_URL")
	if err != nil {
		t.Fatal("failed to unset env variable")
	}

	c, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if c.Port != "8080" {
		t.Errorf("Port = %q, want %q", c.Port, "8080")
	}
	if c.Host != "0.0.0.0" {
		t.Errorf("Host = %q, want %q", c.Host, "0.0.0.0")
	}
	if c.Env != "development" {
		t.Errorf("Env = %q, want %q", c.Env, "development")
	}
	if c.UseSecureCookie != false {
		t.Errorf("UseSecureCookie = %v, want %v", c.UseSecureCookie, false)
	}
	if c.EmailURL != "https://smtp.maileroo.com/api/v2/emails" {
		t.Errorf(
			"EmailURL = %q, want default %q",
			c.EmailURL,
			"https://smtp.maileroo.com/api/v2/emails",
		)
	}

	// Fields without defaults should be empty when unset.
	if c.JwtSecret != "" ||
		c.DatabaseUrl != "" ||
		c.MailerooApiKey != "" ||
		c.EmailReciever != "" ||
		c.EmailSender != "" {
		t.Errorf("expected empty required fields when unset; got: %+v", c)
	}
}

func TestNewConfig_WithEnv(t *testing.T) {
	reset()
	t.Cleanup(reset)

	t.Setenv("PORT", "9090")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("APP_ENV", "production")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("JWT_SECRET", "supersecret")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("EMAIL_MAILEROO_API_KEY", "mailer-key")
	t.Setenv("EMAIL_ACCOUNT", "to@example.com")
	t.Setenv("EMAIL_SENDER_ACCOUNT", "from@example.com")
	t.Setenv("EMAIL_URL", "https://custom.mail/api")

	c, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	tests := map[string]struct {
		got any
		exp any
	}{
		"Port":            {c.Port, "9090"},
		"Host":            {c.Host, "127.0.0.1"},
		"Env":             {c.Env, "production"},
		"UseSecureCookie": {c.UseSecureCookie, true},
		"JwtSecret":       {c.JwtSecret, "supersecret"},
		"DatabaseUrl":     {c.DatabaseUrl, "postgres://user:pass@localhost:5432/db"},
		"MailerooApiKey":  {c.MailerooApiKey, "mailer-key"},
		"EmailReciever":   {c.EmailReciever, "to@example.com"},
		"EmailSender":     {c.EmailSender, "from@example.com"},
		"EmailURL":        {c.EmailURL, "https://custom.mail/api"},
	}

	for name, tt := range tests {
		if tt.got != tt.exp {
			t.Errorf("%s = %#v, want %#v", name, tt.got, tt.exp)
		}
	}
}

func TestNewConfig_Singleton(t *testing.T) {
	reset()
	t.Cleanup(reset)

	t.Setenv("PORT", "1111")
	c1, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	// Change env; since sync.Once is used, NewConfig should NOT re-parse.
	t.Setenv("PORT", "2222")
	c2, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() second call error = %v", err)
	}

	if c1 != c2 {
		t.Fatalf("expected same *Config instance (singleton); got different pointers")
	}
	if c2.Port != "1111" {
		t.Errorf(
			"singleton should not re-parse env; Port = %q, want %q",
			c2.Port,
			"1111",
		)
	}
}

func TestNewConfig_InvalidBool(t *testing.T) {
	reset()
	t.Cleanup(reset)

	t.Setenv("COOKIE_SECURE", "notabool")

	_, err := NewConfig()
	if err == nil {
		t.Fatalf("expected error for invalid boolean in COOKIE_SECURE, got nil")
	}
}
