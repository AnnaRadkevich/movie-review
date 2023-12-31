package movies

import (
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/RadkevichAnn/movie-reviews/internal/modules/stars"

	"github.com/RadkevichAnn/movie-reviews/internal/modules/genres"

	"github.com/RadkevichAnn/movie-reviews/contracts"
	"github.com/RadkevichAnn/movie-reviews/internal/config"
	"github.com/RadkevichAnn/movie-reviews/internal/echox"
	"github.com/RadkevichAnn/movie-reviews/internal/pagination"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig config.PaginationConfig
	reqGroup         singleflight.Group
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
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetOrDeleteMovieByIDRequest](c)
		if err != nil {
			return nil, err
		}
		movie, err := h.service.GetMovieByID(c.Request().Context(), req.ID)
		if err != nil {
			return nil, err
		}
		return movie, err
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllMovies(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetMoviesRequest](c)
		if err != nil {
			return nil, err
		}
		pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
		movies, total, err := h.service.GetAllMoviesPaginated(c.Request().Context(), req.SearchTerm, req.SortByRating, req.StarID, offset, limit)
		if err != nil {
			return nil, err
		}
		return pagination.Response(&req.PaginatedRequest, total, movies), nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
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
