package tests

import (
	"time"

	"github.com/RadkevichAnn/movie-reviews/internal/config"
)

const testPaginationSize = 2

func getConfig(pgConnString string) *config.Config {
	return &config.Config{
		DbUrl: pgConnString,
		Port:  0, // random port
		JWT: config.JwtConfig{
			Secret:           "secret",
			AccessExpiration: time.Minute * 15,
		},
		Admin: config.AdminConfig{
			Username: "admin",
			Password: "&dm1Npa$$",
			Email:    "admin@mail.com",
		},
		Pagination: config.PaginationConfig{
			DefaultSize: testPaginationSize,
			MaxSize:     50,
		},
		Local:    true,
		LogLevel: "error",
	}
}
