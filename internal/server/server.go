package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/modules/movies"

	"github.com/cloudmachinery/movie-reviews/internal/modules/stars"

	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/cloudmachinery/movie-reviews/internal/log"
	"github.com/cloudmachinery/movie-reviews/internal/modules/auth"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/cloudmachinery/movie-reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"gopkg.in/validator.v2"
)

const (
	dbConnectTimeout     = 10 * time.Second
	adminCreationTimeout = 10 * time.Second
)

type Server struct {
	e       *echo.Echo
	cfg     *config.Config
	closers []func() error
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("setup logger: %w", err)
	}
	slog.SetDefault(logger)

	validation.SetupValidators()

	var closers []func() error
	db, err := getDb(ctx, cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}
	closers = append(closers, func() error { db.Close(); return nil })

	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(usersModule.Service, jwtService)
	genresModule := genres.NewModule(db)
	starsModule := stars.NewModule(db, cfg.Pagination)
	moviesModule := movies.NewModule(db, cfg.Pagination, genresModule, starsModule)

	if err = createInitialAdminUser(cfg.Admin, authModule.Service); err != nil {
		return nil, withClosers(closers, fmt.Errorf("create initial admin user: %w", err))
	}

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler

	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	api := e.Group("/api")
	api.Use(jwt.NewAuthMiddleware(cfg.JWT.Secret))
	api.Use(echox.Logger)

	// Auth API routes
	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	// Users API routes
	api.GET("/users/:userId", usersModule.Handler.Get)
	api.GET("/users/username/:username", usersModule.Handler.GetByUserName)
	api.DELETE("/users/:userId", usersModule.Handler.Delete, auth.Self)
	api.PUT("/users/:userId", usersModule.Handler.UpdateBio, auth.Self)
	api.PUT("/users/:userId/role/:role", usersModule.Handler.UpdateRole, auth.Admin)

	// Genres API routes
	api.GET("/genres", genresModule.Handler.GetAllGenres)
	api.GET("/genres/:genreId", genresModule.Handler.GetGenreByID)
	api.POST("/genres", genresModule.Handler.CreateGenre, auth.Editor)
	api.PUT("/genres/:genreId", genresModule.Handler.UpdateGenre, auth.Editor)
	api.DELETE("/genres/:genreId", genresModule.Handler.DeleteGenre, auth.Editor)

	// Stars API routes
	api.GET("/stars", starsModule.Handler.GetAllStars)
	api.GET("/stars/:id", starsModule.Handler.GetStarByID)
	api.POST("/stars", starsModule.Handler.CreateStar, auth.Editor)
	api.PUT("/stars/:id", starsModule.Handler.UpdateStar, auth.Editor)
	api.DELETE("/stars/:id", starsModule.Handler.DeleteStar, auth.Editor)

	// Movies API routes
	api.GET("/movies", moviesModule.Handler.GetAllMovies)
	api.GET("/movies/:id", moviesModule.Handler.GetMovieByID)
	api.POST("/movies", moviesModule.Handler.CreateMovie, auth.Editor)
	api.PUT("/movies/:id", moviesModule.Handler.UpdateMovie, auth.Editor)
	api.DELETE("/movies/:id", moviesModule.Handler.DeleteMovie, auth.Editor)

	return &Server{e: e, cfg: cfg, closers: closers}, nil
}

func (s *Server) Start() error {
	port := s.cfg.Port
	slog.Info("server started", "port", port)
	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Close() error {
	return withClosers(s.closers, nil)
}

func (s *Server) Port() (int, error) {
	listener := s.e.Listener
	if listener == nil {
		return 0, errors.New("server is not started")
	}

	addr := listener.Addr()
	if addr == nil {
		return 0, errors.New("server is not started")
	}

	return addr.(*net.TCPAddr).Port, nil
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

func createInitialAdminUser(cfg config.AdminConfig, authService *auth.Service) error {
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

func withClosers(closers []func() error, err error) error {
	errs := []error{err}

	for i := len(closers) - 1; i >= 0; i-- {
		if err = closers[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
