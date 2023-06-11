package tests

import (
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func StarsApiChecks(t *testing.T, c *client.Client) {
	var lucas, hamill, mcgregor *contracts.Star

	t.Run("stars.CreateStar: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateStarRequest
			addr **contracts.Star
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
}
