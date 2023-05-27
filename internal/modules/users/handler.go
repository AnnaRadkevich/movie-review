package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func (h *Handler) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c echo.Context) error {
	var req DeleteOrGetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	user, err := h.service.Get(c.Request().Context(), req.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c echo.Context) error {
	var req DeleteOrGetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return h.service.Update(c.Request().Context(), req.UserId, req.Bio)
}

type DeleteOrGetRequest struct {
	UserId int `param:"userId"`
}

type UpdateRequest struct {
	UserId int    `param:"userId"`
	Bio    string `json:"bio"`
}
