package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	tokenContextKey = "token"
)

func NewAuthMiddleware(secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey: tokenContextKey,
		SigningKey: []byte(secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &AccessClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return nil
		},
		ContinueOnIgnoredError: true,
	})
}

func GetClaims(c echo.Context) *AccessClaims {
	token := c.Get(tokenContextKey)
	if token == nil {
		return nil
	}
	return token.(*jwt.Token).Claims.(*AccessClaims)
}
