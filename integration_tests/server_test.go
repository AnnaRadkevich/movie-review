package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/server"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	prepareInfrastructure(t, runServer)
}

func runServer(t *testing.T, pgConnString string) {
	cfg := &config.Config{
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
		Local:    true,
		LogLevel: "error",
	}

	srv, err := server.New(context.Background(), cfg)
	require.NoError(t, err)
	defer srv.Close()

	go func() {
		if serr := srv.Start(); serr != http.ErrServerClosed {
			require.NoError(t, serr)
		}
	}()

	var port int
	retry.Run(t, func(r *retry.R) {
		port, err = srv.Port()
		if err != nil {
			require.NoError(r, err)
		}
	})

	tests(t, port, cfg)

	err = srv.Shutdown(context.Background())
	require.NoError(t, err)
}

func tests(t *testing.T, port int, cfg *config.Config) {
	addr := fmt.Sprintf("http://localhost:%d", port)
	c := client.New(addr)

	AuthApiChecks(t, c, cfg)
	UsersApiChecks(t, c, cfg)
	GenresApiChecks(t, c)
	StarsApiChecks(t, c)
}
