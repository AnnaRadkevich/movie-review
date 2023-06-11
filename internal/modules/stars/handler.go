package stars

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/pagination"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
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

func (h *Handler) CreateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateStarRequest](c)
	if err != nil {
		return err
	}
	star := &Star{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		BirthDate:  req.BirthDate,
		BirthPlace: req.BirthPlace,
		DeathDate:  req.DeathDate,
		Bio:        req.Bio,
	}
	err = h.service.CreateStar(c.Request().Context(), star)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, star)
}

func (h *Handler) GetStarByID(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetOrDeleteStarByIDRequest](c)
	if err != nil {
		return err
	}
	star, err := h.service.GetStarByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, star)
}

func (h *Handler) GetAllStars(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetStarsRequest](c)
	if err != nil {
		return err
	}
	pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
	stars, total, err := h.service.GetAllStarsPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pagination.Response(&req.PaginatedRequest, total, stars))
}

func (h *Handler) UpdateStar(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateStarRequest](c)
	if err != nil {
		return err
	}
	star := &Star{
		ID:         req.ID,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		BirthDate:  req.BirthDate,
		BirthPlace: req.BirthPlace,
		DeathDate:  req.DeathDate,
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
