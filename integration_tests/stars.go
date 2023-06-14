package tests

import (
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func StarsApiChecks(t *testing.T, c *client.Client) {
	var lucas, hamill, mcgregor *contracts.StarDetails

	t.Run("stars.CreateStar: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateStarRequest
			addr **contracts.StarDetails
		}{
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "George",
					MiddleName: contracts.Ptr("Walton"),
					LastName:   "Lucas",
					BirthDate:  time.Date(1944, time.May, 14, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Modesto, California6 U.S."),
					Bio:        contracts.Ptr("Famous creator of Star Wars"),
				},
				addr: &lucas,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "Mark",
					MiddleName: contracts.Ptr("Richard"),
					LastName:   "Hamill",
					BirthDate:  time.Date(1951, time.September, 25, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Oakland, California6 U.S."),
				},
				addr: &hamill,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "Ewan",
					MiddleName: contracts.Ptr("Gordon"),
					LastName:   "McGregor",
					BirthDate:  time.Date(1971, time.March, 31, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Perth, Scotland"),
				},
				addr: &mcgregor,
			},
		}

		for _, cc := range cases {

			star, err := c.CreateStar(contracts.NewAuthenticated(cc.req, johnDoeToken))
			require.NoError(t, err)

			*cc.addr = star
			require.NotEmpty(t, star.ID)
			require.NotEmpty(t, star.CreatedAt)
		}
	})
	t.Run("stars.GetStarByID: success", func(t *testing.T) {
		star, err := c.GetStarByID(hamill.ID)
		require.NoError(t, err)
		require.Equal(t, star.ID, hamill.ID)
	})
	t.Run("stars.GetStarByID: not found", func(t *testing.T) {
		notExistingId := 10
		_, err := c.GetStarByID(notExistingId)
		requireNotFoundError(t, err, "star", "id", notExistingId)
	})
	t.Run("stars.GetAllStars: success", func(t *testing.T) {
		req := &contracts.GetStarsRequest{}
		res, err := c.GetStars(req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{&lucas.Star, &hamill.Star}, res.Items)

		req.Page = res.Page + 1
		res, err = c.GetStars(req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 2, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{&mcgregor.Star}, res.Items)
	})
	t.Run("stars.UpdateStar: success", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			ID:         lucas.ID,
			FirstName:  lucas.FirstName,
			MiddleName: lucas.MiddleName,
			LastName:   "LUCAS",
			BirthDate:  lucas.BirthDate,
			DeathDate:  lucas.DeathDate,
			Bio:        contracts.Ptr("UPDATE:Famous creator of Star Wars"),
		}
		err := c.UpdateStar(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		res, err := c.GetStarByID(lucas.ID)
		require.NoError(t, err)
		require.Equal(t, *req.Bio, *res.Bio)
		require.Equal(t, req.LastName, res.LastName)
	})
	t.Run("stars.UpdateStar: not found", func(t *testing.T) {
		notExistingId := 10
		req := &contracts.UpdateStarRequest{
			ID:         notExistingId,
			FirstName:  lucas.FirstName,
			MiddleName: lucas.MiddleName,
			LastName:   lucas.LastName,
			BirthDate:  lucas.BirthDate,
			DeathDate:  lucas.DeathDate,
			Bio:        contracts.Ptr("UPDATE:Famous creator of Star Wars"),
		}
		err := c.UpdateStar(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "star", "id", notExistingId)
	})
	t.Run("stars.DeleteStar: not found", func(t *testing.T) {
		notExistingId := 10
		req := &contracts.GetOrDeleteStarByIDRequest{
			ID: notExistingId,
		}
		err := c.DeleteStar(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "star", "id", notExistingId)
	})
	t.Run("stars.DeleteStar: success", func(t *testing.T) {
		req := &contracts.GetOrDeleteStarByIDRequest{
			ID: lucas.ID,
		}
		err := c.DeleteStar(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err, "star", "id", lucas.ID)

		_, err = c.GetStarByID(lucas.ID)
		requireNotFoundError(t, err, "star", "id", lucas.ID)
	})
}
