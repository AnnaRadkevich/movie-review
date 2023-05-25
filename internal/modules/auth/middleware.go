package auth

import (
	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/labstack/echo/v4"
)

func Self(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		claims := jwt.GetClaims(c)
		if claims.Role == users.AdminRole || claims.Subject == userId {
			return next(c)
		}
		return echo.ErrForbidden
	}
}

func Editor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == users.EditorRole || claims.Role == users.AdminRole {
			return next(c)
		}
		return echo.ErrForbidden
	}
}

func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := jwt.GetClaims(c)
		if claims.Role == users.AdminRole {
			return next(c)
		}
		return echo.ErrForbidden
	}
}
