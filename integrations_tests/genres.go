package integrations_tests

import (
	"net/http"
	"testing"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func GenressApiChecks(t *testing.T, c *client.Client) {
	t.Run("genres.GetAllGenres: empty", func(t *testing.T) {
		genres, err := c.GetAllGenres()
		require.NoError(t, err)
		require.Empty(t, genres)
	})

	var action, drama, spooky *contracts.Genre
	t.Run("genres.CreateGenre: success Action by Admin, Drama and Spooky by John Doe", func(t *testing.T) {
		cases := []struct {
			name  string
			token string
			addr  **contracts.Genre
		}{
			{"Action", adminToken, &action},
			{"Drama", johnDoeToken, &drama},
			{"Spooky", johnDoeToken, &spooky},
		}

		for _, cc := range cases {
			req := &contracts.CreateGenreRequest{
				Name: cc.name,
			}
			g, err := c.CreateGenre(contracts.NewAuthenticated(req, cc.token))
			require.NoError(t, err)

			*cc.addr = g
			require.NotEmpty(t, g.ID)
			require.Equal(t, req.Name, g.Name)
		}
	})
	t.Run("genres.CreateGenre: short name", func(t *testing.T) {
		req := &contracts.CreateGenreRequest{
			Name: "by",
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireBadRequestError(t, err, "Name")
	})
	t.Run("genres:CreateGenre: existing name", func(t *testing.T) {
		req := &contracts.CreateGenreRequest{
			Name: drama.Name,
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireAlreadyExistError(t, err, "genre", "name", req.Name)
	})
	t.Run("genres.GetAllGenres: three genres", func(t *testing.T) {
		genres, err := c.GetAllGenres()
		require.NoError(t, err)
		require.Equal(t, []*contracts.Genre{action, drama, spooky}, genres)
	})
	t.Run("genres.GetGenreById: success", func(t *testing.T) {
		g, err := c.GetGenreById(action.ID)
		require.NoError(t, err)
		require.Equal(t, action, g)
	})
	t.Run("genres.GetGetGenreById: not found", func(t *testing.T) {
		_, err := c.GetGenreById(10)
		requireNotFoundError(t, err, "genre", "id", 10)
	})
	t.Run("genres.UpdateGenre: success", func(t *testing.T) {
		req := &contracts.UpdateGenreRequest{
			GenreId: spooky.ID,
			Name:    "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		spooky = getGenre(t, c, spooky.ID)
		require.Equal(t, req.Name, spooky.Name)
	})

	t.Run("genres.UpdateGenre: not found", func(t *testing.T) {
		nonExistingId := 1000
		req := &contracts.UpdateGenreRequest{
			GenreId: nonExistingId,
			Name:    "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "genre", "id", nonExistingId)
	})
	t.Run("genres.DeleteGenre: success", func(t *testing.T) {
		req := &contracts.DeleteGenreRequest{
			GenreId: spooky.ID,
		}
		err := c.DeleteGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		spooky = getGenre(t, c, spooky.ID)
		require.Nil(t, spooky)
	})
}

func getGenre(t *testing.T, c *client.Client, id int) *contracts.Genre {
	u, err := c.GetGenreById(id)
	if err != nil {
		cerr, ok := err.(*client.Error)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cerr.Code)
		return nil
	}

	return u
}
