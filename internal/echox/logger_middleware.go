package echox

import (
	"github.com/RadkevichAnn/movie-reviews/internal/log"
	"github.com/RadkevichAnn/movie-reviews/internal/modules/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestGroup := slog.Group("request",
			slog.String("method", c.Request().Method),
			slog.String("url", c.Request().URL.String()))
		attrs := []any{requestGroup}

		if claims := jwt.GetClaims(c); claims != nil {
			requesterGroup := slog.Group("requester",
				slog.Int("id", claims.UserID),
				slog.String("role", claims.Role))
			attrs = append(attrs, requesterGroup)
		}
		logger := slog.Default().With(attrs...)
		ctx := log.WithLogger(c.Request().Context(), logger)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}
