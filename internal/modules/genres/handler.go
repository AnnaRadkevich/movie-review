package genres

import (
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  *Service
	reqGroup singleflight.Group
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllGenres(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		genres, err := h.service.GetAllGenres(c.Request().Context())
		if err != nil {
			return nil, err
		}

		return genres, nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetGenreByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetGenreRequest](c)
	if err != nil {
		return err
	}
	genre, err := h.service.GetGenreById(c.Request().Context(), req.GenreId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, genre)
}

func (h *Handler) CreateGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateGenreRequest](c)
	if err != nil {
		return err
	}
	genre, err := h.service.CreateGenre(c.Request().Context(), req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, genre)
}

func (h *Handler) UpdateGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	return h.service.UpdateGenre(c.Request().Context(), req.GenreId, req.Name)
}

func (h *Handler) DeleteGenre(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteGenreRequest](c)
	if err != nil {
		return err
	}
	return h.service.DeleteGenre(c.Request().Context(), req.GenreId)
}
