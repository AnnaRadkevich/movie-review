package users

import (
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteOrGetRequest](c)
	if err != nil {
		return err
	}
	user, err := h.service.Get(c.Request().Context(), req.UserId)
	if err != nil {
		return apperrors.BadRequest(err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteOrGetRequest](c)
	if err != nil {
		return err
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h *Handler) UpdateBio(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateRequest](c)
	if err != nil {
		return err
	}
	return h.service.UpdateBio(c.Request().Context(), req.UserId, *req.Bio)
}

func (h *Handler) UpdateRole(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateRoleRequest](c)
	if err != nil {
		return err
	}
	return h.service.UpdateRole(c.Request().Context(), req.UserID, req.Role)
}

func (h *Handler) GetByUserName(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetUserByUsernameRequest](c)
	if err != nil {
		return err
	}

	user, err := h.service.GetExistingUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
