package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	DbUrl string    `env:"DB_URL"`
	Port  int       `env:"PORT" envDefault:"8080"`
	JWT   JwtConfig `envPrefix:"JWT_"`
}

type JwtConfig struct {
	Secret           string        `env:"SECRET"'`
	AccessExpiration time.Duration `env:"ACCESS_EXPIRATION" envDefault:"15m"`
}

func NewConfig() (*Config, error) {
	var c Config
	err := env.Parse(&c)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &c, nil
}
