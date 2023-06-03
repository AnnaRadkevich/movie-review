package auth

import (
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

var errForbidden = apperrors.Forbidden("insufficient permissions")

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		userId := c.Param("userId")
		if claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}
		return errForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		if claims.Role == users.EditorRole || claims.Role == users.AdminRole {
			return next(c)
		}
		return errForbidden
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims == nil {
			return apperrors.Unauthorized("invalid or missing token")
		}
		if claims.Role == users.AdminRole {
			return next(c)
		}
		return errForbidden
	}
}
