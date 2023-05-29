package auth

import (
	"github.com/cloudmachinery/movie-reviews/internal/echox"
	"net/http"

	"gopkg.in/validator.v2"

	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(c echo.Context) error {
	req, err := echox.BindAndValidate[RegisterRequest](c)
	if err != nil {
		return err
	}
	if err := validator.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user := &users.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.authService.Register(c.Request().Context(), user, req.Password); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	req, err := echox.BindAndValidate[LoginRequest](c)
	if err != nil {
		return err
	}
	if err := validator.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	accessToken, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}

	response := LoginResponse{
		AccessToken: accessToken,
	}
	return c.JSON(http.StatusOK, response)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"min=5,max=16" `
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
