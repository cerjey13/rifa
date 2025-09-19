package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server    ServerOpts
	Database  DatabaseOpts
	Service   ServiceOpts
	Collector CollectorOpts
}

type ServerOpts struct {
	Port     string `env:"PORT" envDefault:"8080"`
	Host     string `env:"HOST" envDefault:"0.0.0.0"`
	Env      string `env:"APP_ENV" envDefault:"development"`
	TimeOuts struct {
		Write      time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
		Read       time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
		ReadHeader time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
		Idle       time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	}
}

type DatabaseOpts struct {
	DatabaseUrl string `env:"DATABASE_URL"`
}
type ServiceOpts struct {
	UseSecureCookie bool `env:"COOKIE_SECURE" envDefault:"false"`
	JwtOpts         JwtOpts
	Email           EmailOpts
}

type JwtOpts struct {
	JwtSecret    string `env:"JWT_SECRET"`
	JwtExpiresAt int    `env:"JWT_EXPIRES_AT" envDefault:"168"`
}

type EmailOpts struct {
	MailerooApiKey string `env:"EMAIL_MAILEROO_API_KEY"`
	EmailReciever  string `env:"EMAIL_ACCOUNT"`
	EmailSender    string `env:"EMAIL_SENDER_ACCOUNT"`
	EmailURL       string `env:"EMAIL_URL" envDefault:"https://smtp.maileroo.com/api/v2/emails"`
}

type CollectorOpts struct {
	CollectorEnv             string `env:"APP_ENV" envDefault:"development"`
	CollectorExporter        string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	CollectorExporterHeaders string `env:"OTEL_EXPORTER_OTLP_HEADERS"`
}

var (
	once sync.Once
	cfg  *Config = &Config{}
)

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() (*Config, error) {
	var err error
	once.Do(func() {
		err = env.Parse(cfg)
	})
	if err != nil {
		return nil, err
	}

	return cfg, err
}
