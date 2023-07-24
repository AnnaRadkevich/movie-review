package tests

import (
	"net/http"
	"testing"

	"github.com/RadkevichAnn/movie-reviews/client"
	"github.com/RadkevichAnn/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

var (
	Action *contracts.Genre
	Drama  *contracts.Genre
	Spooky *contracts.Genre
)

func GenresApiChecks(t *testing.T, c *client.Client) {
	t.Run("genres.GetAllGenres: empty", func(t *testing.T) {
		genres, err := c.GetAllGenres()
		require.NoError(t, err)
		require.Empty(t, genres)
	})

	t.Run("genres.CreateGenre: success Action by Admin, Drama and Spooky by John Doe", func(t *testing.T) {
		cases := []struct {
			name  string
			token string
			addr  **contracts.Genre
		}{
			{"Action", adminToken, &Action},
			{"Drama", johnDoeToken, &Drama},
			{"Spooky", johnDoeToken, &Spooky},
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
			Name: Drama.Name,
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireAlreadyExistError(t, err, "genre", "name", req.Name)
	})
	t.Run("genres.GetAllGenres: three genres", func(t *testing.T) {
		genres, err := c.GetAllGenres()
		require.NoError(t, err)
		require.Equal(t, []*contracts.Genre{Action, Drama, Spooky}, genres)
	})
	t.Run("genres.GetGenreByID: success", func(t *testing.T) {
		g, err := c.GetGenreById(Action.ID)
		require.NoError(t, err)
		require.Equal(t, Action, g)
	})
	t.Run("genres.GetGetGenreById: not found", func(t *testing.T) {
		_, err := c.GetGenreById(10)
		requireNotFoundError(t, err, "genre", "id", 10)
	})
	t.Run("genres.UpdateGenre: success", func(t *testing.T) {
		req := &contracts.UpdateGenreRequest{
			GenreId: Spooky.ID,
			Name:    "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		Spooky = getGenre(t, c, Spooky.ID)
		require.Equal(t, req.Name, Spooky.Name)
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
			GenreId: Spooky.ID,
		}
		err := c.DeleteGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		Spooky = getGenre(t, c, Spooky.ID)
		require.Nil(t, Spooky)
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
