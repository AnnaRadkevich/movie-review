package users

import (
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c echo.Context) error {
	req, err := echox.BindAndValidate[DeleteOrGetRequest](c)
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
	req, err := echox.BindAndValidate[DeleteOrGetRequest](c)
	if err != nil {
		return err
	}
	return h.service.Delete(c.Request().Context(), req.UserId)
}

func (h *Handler) UpdateBio(c echo.Context) error {
	req, err := echox.BindAndValidate[UpdateRequest](c)
	if err != nil {
		return err
	}
	return h.service.UpdateBio(c.Request().Context(), req.UserId, req.Bio)
}
func (h *Handler) UpdateRole(c echo.Context) error {
	req, err := echox.BindAndValidate[UpdateRoleRequest](c)
	if err != nil {
		return err
	}
	return h.service.UpdateRole(c.Request().Context(), req.UserID, req.Role)
}

type DeleteOrGetRequest struct {
	UserId int `param:"userId"`
}

type UpdateRequest struct {
	UserId int    `param:"userId"`
	Bio    string `json:"bio"`
}

type UpdateRoleRequest struct {
	UserID int    `param:"userId" validate:"nonzero"`
	Role   string `param:"role" validate:"role"`
}
