package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/cloudmachinery/movie-reviews/internal/log"
	"github.com/cloudmachinery/movie-reviews/internal/validation"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"gopkg.in/validator.v2"

	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	dbConnectTimeout     = 10 * time.Second
	gracefulTimeout      = 10 * time.Second
	adminCreationTimeout = 10 * time.Second
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	validation.SetupValidators()

	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	failOnError(err, "setup logger")
	slog.SetDefault(logger)

	slog.Info("started", "config", cfg)

	db, err := getDb(context.Background(), cfg.DbUrl)
	failOnError(err, "connect to db")
	defer db.Close()

	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)

	err = createAdmin(cfg.Admin, authModule.Service)
	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler

	e.Use(middleware.Recover())
	api := e.Group("/api")
	api.Use(jwt.NewAuthMiddleware(cfg.JWT.Secret))

	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	api.GET("/users/:userId", usersModule.Handler.Get)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.UpdateBio, auth.Self)
	api.PUT("/users/:userID/role/:role", usersModule.Handler.UpdateRole, auth.Admin)

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel

		ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "err", err)
		}
	}()
	if err = e.Start(fmt.Sprintf(":%d", cfg.Port)); err != http.ErrServerClosed {
		slog.Error("server shutdown caused by", "err", err)
	}
}

func createAdmin(cfg config.AdminConfig, authService *auth.Service) error {
	if !cfg.IsSet() {
		return nil
	}
	if err := validator.Validate(cfg); err != nil {
		return fmt.Errorf("validate admin config: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), adminCreationTimeout)
	defer cancel()

	err := authService.Register(ctx, &users.User{
		Username: cfg.Username,
		Email:    cfg.Email,
		Role:     users.AdminRole,
	}, cfg.Password)

	switch {
	case apperrors.Is(err, apperrors.InternalCode):
		return fmt.Errorf("register admin :%w", err)
	case err != nil:
		return nil
	default:
		slog.Info("admin user created", "username", cfg.Username, "email", cfg.Email)
		return nil

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
		slog.Error(message, "err", err)
		os.Exit(1)
	}
}
