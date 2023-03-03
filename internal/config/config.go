package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN"` // Database connection string for DB-style storage
	APIKey      string `env:"API_KEY"`
}

func Init() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "secret key")
	flag.StringVar(&cfg.APIKey, "a", cfg.APIKey, "api key")

	flag.Parse()

	return cfg, nil
}
