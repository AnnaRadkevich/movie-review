package integrations_tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/cloudmachinery/movie-reviews/internal/server"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	prepareInfrastructure(t, runServer)
}

func runServer(t *testing.T, pgConnString string) {
	cfg := &config.Config{
		DbUrl: pgConnString,
		Port:  0, // random port
		JWT: config.JwtConfig{
			Secret:           "secret",
			AccessExpiration: time.Minute * 15,
		},
		Admin: config.AdminConfig{
			Username: "admin",
			Password: "&dm1Npa$$",
			Email:    "admin@mail.com",
		},
		Local:    true,
		LogLevel: "error",
	}

	srv, err := server.New(context.Background(), cfg)
	require.NoError(t, err)
	defer srv.Close()

	go func() {
		if serr := srv.Start(); serr != http.ErrServerClosed {
			require.NoError(t, serr)
		}
	}()

	var port int
	retry.Run(t, func(r *retry.R) {
		port, err = srv.Port()
		if err != nil {
			require.NoError(r, err)
		}
	})

	tests(t, port, cfg)

	err = srv.Shutdown(context.Background())
	require.NoError(t, err)
}

func tests(t *testing.T, port int, cfg *config.Config) {
	addr := fmt.Sprintf("http://localhost:%d", port)
	c := client.New(addr)

	t.Run("users.GetUserByUserName: admin", func(t *testing.T) {
		u, err := c.GetUserByUserName(cfg.Admin.Username)
		require.NoError(t, err)

		require.Equal(t, cfg.Admin.Username, u.Username)
		require.Equal(t, cfg.Admin.Email, u.Email)
		require.Equal(t, users.AdminRole, u.Role)
	})

	t.Run("users.GetUserByUserName: not found", func(t *testing.T) {
		_, err := c.GetUserByUserName("notfound")
		requireNotFoundError(t, err, "user", "username", "notfound")
	})

	var (
		johnDoe     *contracts.User
		johnDoePass = standardPassword
	)
	t.Run("auth.Register: success", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johndoe",
			Email:    "johndoe@example.com",
			Password: johnDoePass,
		}
		u, err := c.RegisterUser(req)
		require.NoError(t, err)
		johnDoe = u

		require.Equal(t, req.Username, u.Username)
		require.Equal(t, req.Email, u.Email)
		require.Equal(t, users.UserRole, u.Role)
	})

	t.Run("auth.Register: short username", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "joh",
			Email:    "joh@example.com",
			Password: standardPassword,
		}
		_, err := c.RegisterUser(req)
		requireBadRequestError(t, err, "Username")
	})

	var johnDoeToken string
	t.Run("auth.Login: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    johnDoe.Email,
			Password: johnDoePass,
		}
		res, err := c.LoginUser(req)
		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		johnDoeToken = res.AccessToken
	})

	t.Run("users.UpdateUser: success", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)
	})

	t.Run("users.UpdateUser: non-authenticated", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("users.UpdateUser: another user", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateRequest{
			UserId: johnDoe.ID + 1,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})
	var adminToken string
	t.Run("auth.Login: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    cfg.Admin.Email,
			Password: cfg.Admin.Password,
		}
		res, err := c.LoginUser(req)
		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		adminToken = res.AccessToken
	})
	t.Run("users.UpdateRole: success", func(t *testing.T) {
		req := &contracts.UpdateRoleRequest{
			UserID: johnDoe.ID,
			Role:   users.AdminRole,
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)
	})
	t.Run("users.Delete: another user", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{UserId: johnDoe.ID + 1}

		err := c.Delete(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})
	t.Run("users.Delete: non-authenticated", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{UserId: johnDoe.ID}
		err := c.Delete(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})
	t.Run("users.Delete: success", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{UserId: johnDoe.ID}
		err := c.Delete(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)
	})
}

const standardPassword = "secuR3P@ss"

func requireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireApiError(t, err, http.StatusNotFound, msg)
}

func requireUnauthorizedError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusUnauthorized, msg)
}

func requireForbiddenError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusForbidden, msg)
}

func requireBadRequestError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusBadRequest, msg)
}

func requireApiError(t *testing.T, err error, statusCode int, msg string) {
	cerr, ok := err.(*client.Error)
	require.True(t, ok, "expected client.Error")
	require.Equal(t, statusCode, cerr.Code)
	require.Contains(t, cerr.Message, msg)
}
