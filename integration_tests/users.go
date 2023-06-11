package tests

import (
	"net/http"
	"testing"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"github.com/stretchr/testify/require"
)

func UsersApiChecks(t *testing.T, c *client.Client, cfg *config.Config) {
	t.Run("users.GetUserByUserName: admin", func(t *testing.T) {
		u, err := c.GetUserByUserName(cfg.Admin.Username)
		require.NoError(t, err)

		require.Equal(t, cfg.Admin.Username, u.Username)
		require.Equal(t, cfg.Admin.Email, u.Email)
		require.Equal(t, users.AdminRole, u.Role)
	})

	t.Run("users.GetUserByUserName: not found", func(t *testing.T) {
		_, err := c.GetUserByUserName("not found")
		requireNotFoundError(t, err, "user", "username", "not found")
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
	t.Run("users.UpdateUserRole: John Doe to editor", func(t *testing.T) {
		req := &contracts.UpdateRoleRequest{
			UserID: johnDoe.ID,
			Role:   users.EditorRole,
		}
		err := c.UpdateUserRole(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)

		johnDoe = getUser(t, c, johnDoe.ID)
		require.Equal(t, users.EditorRole, johnDoe.Role)

		// Have to re-login to become an editor
		johnDoeToken = login(t, c, johnDoe.Email, johnDoePass)
	})
	randomUser := RegisterRandomUser(t, c)
	t.Run("users.DeleteUser: another user", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{
			UserId: randomUser.ID,
		}
		err := c.Delete(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")

		randomUser = getUser(t, c, randomUser.ID)
		require.NotNil(t, randomUser)
	})

	t.Run("users.Delete: non-authenticated", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{
			UserId: randomUser.ID,
		}
		err := c.Delete(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})
	t.Run("users.DeleteUser: by admin", func(t *testing.T) {
		req := &contracts.DeleteOrGetRequest{
			UserId: randomUser.ID,
		}
		err := c.Delete(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)

		randomUser = getUser(t, c, randomUser.ID)
		require.Nil(t, randomUser)
	})
}

func getUser(t *testing.T, c *client.Client, id int) *contracts.User {
	u, err := c.GetUser(id)
	if err != nil {
		cerr, ok := err.(*client.Error)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cerr.Code)
		return nil
	}

	return u
}
