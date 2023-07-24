package stars

import (
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/RadkevichAnn/movie-reviews/internal/config"
	"github.com/RadkevichAnn/movie-reviews/internal/pagination"

	"github.com/RadkevichAnn/movie-reviews/contracts"
	"github.com/RadkevichAnn/movie-reviews/internal/echox"
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

func (h *Handler) CreateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateStarRequest](c)
	if err != nil {
		return err
	}
	star := &StarDetails{
		Star: Star{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: req.BirthDate,
			DeathDate: req.DeathDate,
		},
		MiddleName: req.MiddleName,
		BirthPlace: req.BirthPlace,
		Bio:        req.Bio,
	}
	err = h.service.CreateStar(c.Request().Context(), star)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, star)
}

func (h *Handler) GetStarByID(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetOrDeleteStarByIDRequest](c)
		if err != nil {
			return nil, err
		}
		star, err := h.service.GetStarByID(c.Request().Context(), req.ID)
		if err != nil {
			return nil, err
		}
		return star, nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllStars(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetStarsRequest](c)
		if err != nil {
			return nil, err
		}
		pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
		stars, total, err := h.service.GetAllStarsPaginated(c.Request().Context(), req.MovieID, offset, limit)
		if err != nil {
			return nil, err
		}
		return pagination.Response(&req.PaginatedRequest, total, stars), nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateStarRequest](c)
	if err != nil {
		return err
	}
	star := &StarDetails{
		Star: Star{
			ID:        req.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			BirthDate: req.BirthDate,
			DeathDate: req.DeathDate,
		},
		MiddleName: req.MiddleName,
		BirthPlace: req.BirthPlace,
		Bio:        req.Bio,
	}
	if err = h.service.UpdateStar(c.Request().Context(), star); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetOrDeleteStarByIDRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteStar(c.Request().Context(), req.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
