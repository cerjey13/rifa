package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	Env         string `env:"APP_ENV" envDefault:"development"`
	JwtSecret   string `env:"JWT_SECRET"`
	DatabaseUrl string `env:"DATABASE_URL"`
}

var once sync.Once

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() (*Config, error) {
	var (
		err error
		cfg *Config = &Config{}
	)
	once.Do(func() {
		err = env.Parse(cfg)
	})
	if err != nil {
		return nil, err
	}

	return cfg, err
}
