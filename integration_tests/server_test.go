package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/RadkevichAnn/movie-reviews/client"
	"github.com/RadkevichAnn/movie-reviews/internal/config"
	"github.com/RadkevichAnn/movie-reviews/internal/server"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	prepareInfrastructure(t, runServer)
}

func runServer(t *testing.T, pgConnString string) {
	cfg := getConfig(pgConnString)

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
	moviesAPIChecks(t, c)
	reviewsAPIChecks(t, c)
}
