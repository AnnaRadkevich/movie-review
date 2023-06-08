package stars

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
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
	req, err := echox.BindAndValidate[contracts.GetStarByIDRequest](c)
	if err != nil {
		return err
	}
	star, err := h.service.GetStarByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, star)
}
