package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/server"
	"golang.org/x/exp/slog"

	"github.com/cloudmachinery/movie-reviews/internal/config"
)

const gracefulTimeout = 10 * time.Second

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	srv, err := server.New(context.Background(), cfg)
	failOnError(err, "create server")

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "err", err)
		}
	}()
	if err = srv.Start(); err != http.ErrServerClosed {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func failOnError(err error, message string) {
	if err != nil {
		slog.Error(message, "err", err)
		os.Exit(1)
	}
}
