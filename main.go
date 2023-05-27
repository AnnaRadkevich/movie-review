package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/validation"

	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	dbConnectTimeout = 10 * time.Second
	gracefulTimeout  = 10 * time.Second
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	validation.SetupValidators()

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")
	defer db.Close()

	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)
	authMiddleware := jwt.NewAuthMiddleware(cfg.JWT.Secret)

	e := echo.New()
	api := e.Group("/api")

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	api.GET("/users/:userId", usersModule.Handler.Get)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, authMiddleware, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.Update, authMiddleware, auth.Self)

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown server: %s", err)
		}
	}()
	if err = e.Start(fmt.Sprintf(":%d", cfg.Port)); err != http.ErrServerClosed {
		log.Printf("server shutdown caused by:%s", err)
	}
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
