package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const dbConnectTimeout = 10 * time.Second

func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")

	err = db.Ping(context.Background())
	failOnError(err, "ping db")

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := e.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutdown server: %s", err)
		}
	}()
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	failOnError(err, "get config")
}

func getDb(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	return db, nil
}
func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
	}
}
