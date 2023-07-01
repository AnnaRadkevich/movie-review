package movies

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/internal/modules/stars"

	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/cloudmachinery/movie-reviews/internal/pagination"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig config.PaginationConfig
}

func NewHandler(service *Service, paginationConfig config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) CreateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateMovieRequest](c)
	if err != nil {
		return err
	}
	movie := &MovieDetails{
		Movie: Movie{
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
		},
		Description: req.Description,
	}
	for _, genreID := range req.Genres {
		movie.Genres = append(movie.Genres, &genres.Genre{ID: genreID})
	}
	for _, creditID := range req.Cast {
		movie.Cast = append(movie.Cast, &stars.MovieCredit{
			Star: stars.Star{
				ID: creditID.StarID,
			},
			Role:    creditID.Role,
			Details: creditID.Details,
		})
	}
	err = h.service.CreateMovie(c.Request().Context(), movie)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, movie)
}

func (h *Handler) GetMovieByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetOrDeleteMovieByIDRequest](c)
	if err != nil {
		return err
	}
	movie, err := h.service.GetMovieByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movie)
}

func (h *Handler) GetAllMovies(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
	if err != nil {
		return err
	}
	pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
	movies, total, err := h.service.GetAllMoviesPaginated(c.Request().Context(), req.SearchTerm, req.StarID, offset, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pagination.Response(&req.PaginatedRequest, total, movies))
}

func (h *Handler) UpdateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateMovieRequest](c)
	if err != nil {
		return err
	}
	movie := &MovieDetails{
		Movie: Movie{
			ID:          req.ID,
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
		},
		Description: req.Description,
	}
	for _, genreID := range req.Genres {
		movie.Genres = append(movie.Genres, &genres.Genre{ID: genreID})
	}
	for _, creditID := range req.Cast {
		movie.Cast = append(movie.Cast, &stars.MovieCredit{
			Star: stars.Star{
				ID: creditID.StarID,
			},
			Role:    creditID.Role,
			Details: creditID.Details,
		})
	}
	if err = h.service.UpdateMovie(c.Request().Context(), movie); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetOrDeleteMovieByIDRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteMovie(c.Request().Context(), req.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
