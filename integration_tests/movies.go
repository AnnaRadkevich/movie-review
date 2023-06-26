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
					Genres:      []int{Action.ID, Drama.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID: mcgregor.ID,
							Role:   "director",
						},
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char1, char2"),
						},
					},
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
					Genres:      []int{Action.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char3"),
						},
					},
				},
				addr: &harryPotter,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title: "The Lord of the Rings. The Fellowship of the Ring",
					Description: "The Lord of the Rings is a series of three epic fantasy adventure films directed by Peter Jackson," +
						" based on the novel The Lord of the Rings by J. R. R. Tolkien",
					ReleaseDate: time.Date(2001, time.December, 10, 0, 0, 0, 0, time.UTC),
					Genres:      []int{Drama.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char4"),
						},
					},
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
			require.NotEmpty(t, movie.Genres)
			require.Equal(t, len(cc.req.Genres), len(movie.Genres))
			require.NotEmpty(t, movie.Cast)
			require.Equal(t, len(cc.req.Cast), len(movie.Cast))
		}
	})

	t.Run("movies.GetMovieByID: success", func(t *testing.T) {
		movie, err := c.GetMovieByID(harryPotter.ID)
		require.NoError(t, err)
		require.Equal(t, movie.ID, harryPotter.ID)
		require.Equal(t, len(harryPotter.Genres), len(movie.Genres))
		for i, genre := range harryPotter.Genres {
			require.Equal(t, *genre, *movie.Genres[i])
		}
		require.Equal(t, len(harryPotter.Cast), len(movie.Cast))
		for i, cast := range harryPotter.Cast {
			require.Equal(t, *cast, *movie.Cast[i])
		}
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
	t.Run("stars.GetAllStars: by movieId success", func(t *testing.T) {
		req := contracts.GetStarsRequest{
			MovieID: contracts.Ptr(lordOfTheRing.ID),
		}
		res, err := c.GetStars(&req)
		require.NoError(t, err)
		require.Equal(t, len(lordOfTheRing.Cast), res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{&hamill.Star}, res.Items)
	})
	t.Run("movies.UpdateMovie: success", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          harryPotter.ID,
			Title:       "Harry Potter and the Philosopher's Stone",
			ReleaseDate: harryPotter.ReleaseDate,
			Description: "is a 2001 fantasy film directed by Chris Columbus and produced by David Heyman," +
				" from a screenplay by Steve Kloves, based on the 1997 novel of the same name by J. K. Rowling." +
				" It is the first installment in the Harry Potter film series. ",
			Genres: []int{Action.ID, Drama.ID},
			Cast: []*contracts.MovieCreditInfo{
				{
					StarID: mcgregor.ID,
					Role:   "producer",
				},
				{
					StarID: hamill.ID,
					Role:   "director",
				},
			},
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		err = c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireVersionMismatchError(t, err, "movie", "id", req.ID, req.Version)

		res, err := c.GetMovieByID(harryPotter.ID)
		require.NoError(t, err)
		require.Equal(t, req.Title, res.Title)
		require.Equal(t, []*contracts.Genre{Action, Drama}, res.Genres)
		require.Equal(t, 1, res.Version)
		for i, credit := range req.Cast {
			require.Equal(t, credit.StarID, res.Cast[i].Star.ID)
			require.NotNil(t, res.Cast[i].Star.FirstName)
			require.NotNil(t, res.Cast[i].Star.LastName)
			require.Equal(t, credit.Role, req.Cast[i].Role)
			require.Equal(t, credit.Details, res.Cast[i].Details)

		}
	})
	t.Run("movies.UpdateMovie: not found", func(t *testing.T) {
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
	t.Run("movies.DeleteMovie: not found", func(t *testing.T) {
		notExistingId := 10
		req := &contracts.GetOrDeleteMovieByIDRequest{
			ID: notExistingId,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "movie", "id", notExistingId)
	})
	t.Run("movies.DeleteMovie: success", func(t *testing.T) {
		req := &contracts.GetOrDeleteMovieByIDRequest{
			ID: harryPotter.ID,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err, "movie", "id", harryPotter.ID)

		_, err = c.GetMovieByID(harryPotter.ID)
		requireNotFoundError(t, err, "movie", "id", harryPotter.ID)
	})
}
