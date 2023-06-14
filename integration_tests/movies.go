package tests

import (
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func moviesAPIChecks(t *testing.T, c *client.Client) {
	var starWars, harryPotter, lordOfTheRing *contracts.MovieDetails
	t.Run("movies.CreateMovie: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateMovieRequest
			addr **contracts.MovieDetails
		}{
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Star Wars",
					Description: "Star Wars is an American epic space opera",
					ReleaseDate: time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC),
				},
				addr: &starWars,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title: "Harry Poster and the Philosopher's Stone",
					Description: "is a 2001 fantasy film directed by Chris Columbus and produced by David Heyman," +
						" from a screenplay by Steve Kloves, based on the 1997 novel of the same name by J. K. Rowling." +
						" It is the first installment in the Harry Potter film series. ",
					ReleaseDate: time.Date(2001, time.November, 4, 0, 0, 0, 0, time.UTC),
				},
				addr: &harryPotter,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title: "The Lord of the Rings. The Fellowship of the Ring",
					Description: "The Lord of the Rings is a series of three epic fantasy adventure films directed by Peter Jackson," +
						" based on the novel The Lord of the Rings by J. R. R. Tolkien",
					ReleaseDate: time.Date(2001, time.December, 10, 0, 0, 0, 0, time.UTC),
				},
				addr: &lordOfTheRing,
			},
		}

		for _, cc := range cases {

			movie, err := c.CreateMovie(contracts.NewAuthenticated(cc.req, johnDoeToken))
			require.NoError(t, err)

			*cc.addr = movie
			require.NotEmpty(t, movie.ID)
			require.NotEmpty(t, movie.CreatedAt)
		}
	})

	t.Run("movies.GetMovieByID: success", func(t *testing.T) {
		movie, err := c.GetMovieByID(harryPotter.ID)
		require.NoError(t, err)
		require.Equal(t, movie.ID, harryPotter.ID)
	})
	t.Run("movies.GetMovieByID: not found", func(t *testing.T) {
		notExistingId := 10
		_, err := c.GetMovieByID(notExistingId)
		requireNotFoundError(t, err, "movie", "id", notExistingId)
	})

	t.Run("movies.GetAllmovies: success", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{}
		res, err := c.GetMovies(req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Movie{&starWars.Movie, &harryPotter.Movie}, res.Items)

		req.Page = res.Page + 1
		res, err = c.GetMovies(req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 2, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Movie{&lordOfTheRing.Movie}, res.Items)
	})
	t.Run("stars.UpdateMovie: success", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          harryPotter.ID,
			Title:       "Harry Potter and the Philosopher's Stone",
			ReleaseDate: harryPotter.ReleaseDate,
			Description: "is a 2001 fantasy film directed by Chris Columbus and produced by David Heyman," +
				" from a screenplay by Steve Kloves, based on the 1997 novel of the same name by J. K. Rowling." +
				" It is the first installment in the Harry Potter film series. ",
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		res, err := c.GetMovieByID(harryPotter.ID)
		require.NoError(t, err)
		require.Equal(t, req.Title, res.Title)
	})
	t.Run("stars.UpdateMovie: not found", func(t *testing.T) {
		notExistingId := 10
		req := &contracts.UpdateMovieRequest{
			ID:          notExistingId,
			Title:       " ",
			Description: " ",
			ReleaseDate: time.Date(2000, time.May, 1, 0, 0, 0, 0, time.UTC),
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "movie", "id", notExistingId)
	})
	t.Run("stars.DeleteMovie: not found", func(t *testing.T) {
		notExistingId := 10
		req := &contracts.GetOrDeleteMovieByIDRequest{
			ID: notExistingId,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "movie", "id", notExistingId)
	})
	t.Run("stars.DeleteMOvie: success", func(t *testing.T) {
		req := &contracts.GetOrDeleteMovieByIDRequest{
			ID: harryPotter.ID,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err, "movie", "id", harryPotter.ID)

		_, err = c.GetMovieByID(harryPotter.ID)
		requireNotFoundError(t, err, "movie", "id", harryPotter.ID)
	})
}
