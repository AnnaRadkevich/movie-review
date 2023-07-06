package tests

import (
	"net/http"
	"testing"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func reviewsAPIChecks(t *testing.T, c *client.Client) {
	reviewer1 := RegisterRandomUser(t, c)
	reviewer2 := RegisterRandomUser(t, c)
	reviewer1Token := login(t, c, reviewer1.Email, standardPassword)
	reviewer2Token := login(t, c, reviewer2.Email, standardPassword)

	var review1, review2, review3 *contracts.Review

	t.Run("reviews.CreateReview: success", func(t *testing.T) {
		cases := []struct {
			req   *contracts.CreateReviewRequest
			token string
			addr  **contracts.Review
		}{
			{
				req: &contracts.CreateReviewRequest{
					MovieID: starWars.ID,
					UserID:  reviewer1.ID,
					Rating:  9,
					Title:   "Legendary piece of cinema",
					Content: "I love the original Star Wars films! They're a magical experience with great music and " +
						"sounds. They were made amazingly for their time. Some parts can be boring, but overall " +
						"they're glorious. I didn't understand the hype until a few years ago, but now I'm happy with" +
						"all the films, including the new ones.",
				},
				token: reviewer1Token,
				addr:  &review1,
			},
			{
				req: &contracts.CreateReviewRequest{
					MovieID: starWars.ID,
					UserID:  reviewer2.ID,
					Rating:  8,
					Title:   "A long time ago in a decade without CGI...",
					Content: "A timeless classic with impressive practical effects, despite outdated CGI. A must-watch " +
						"for fans of the franchise and a testament to its enduring greatness.",
				},
				token: reviewer2Token,
				addr:  &review2,
			},
			{
				req: &contracts.CreateReviewRequest{
					MovieID: lordOfTheRing.ID,
					UserID:  reviewer1.ID,
					Rating:  10,
					Title:   "Legendary piece of cinema",
					Content: "Jackson has put together an amazing cast. Particularly pleased with Ian McKellen. " +
						"His game, despite the absence of strong and vivid emotions, fascinates. " +
						"Forgive me, Viggo fans, but I think he's overdoing it a bit. " +
						"He turned out to be a somewhat vain hero. It's probably the fault of the writers, though.",
				},
				token: reviewer1Token,
				addr:  &review3,
			},
		}

		for _, cc := range cases {
			review, err := c.CreateReview(contracts.NewAuthenticated(cc.req, cc.token))
			require.NoError(t, err)

			*cc.addr = review
			require.NotEmpty(t, review.MovieID)
			require.NotEmpty(t, review.UserID)
			require.NotEmpty(t, review.Title)
			require.NotEmpty(t, review.Content)
		}
	})
	t.Run("reviews.GetReviewByID: success", func(t *testing.T) {
		for _, review := range []*contracts.Review{review1, review2, review3} {
			r, err := c.GetReviewByID(review.ID)
			require.NoError(t, err)

			require.Equal(t, review, r)
		}
	})
	t.Run("reviews.GetReviewByID: not found", func(t *testing.T) {
		notExistingId := 10
		_, err := c.GetReviewByID(notExistingId)
		requireNotFoundError(t, err, "review", "id", notExistingId)
	})
	t.Run("reviews.GetAllReviewsPaginated: success", func(t *testing.T) {
		cases := []struct {
			req *contracts.GetReviewsRequest
			exp []*contracts.Review
		}{
			{
				req: &contracts.GetReviewsRequest{
					MovieID: contracts.Ptr(starWars.ID),
				},
				exp: []*contracts.Review{review1, review2},
			},
			{
				req: &contracts.GetReviewsRequest{
					MovieID: contracts.Ptr(lordOfTheRing.ID),
				},
				exp: []*contracts.Review{review3},
			},
			{
				req: &contracts.GetReviewsRequest{
					UserID: contracts.Ptr(reviewer1.ID),
				},
				exp: []*contracts.Review{review1, review3},
			},
		}
		for _, cc := range cases {
			res, err := c.GetAllReviewsPaginated(cc.req)
			require.NoError(t, err)

			require.Equal(t, cc.exp, res.Items)
		}
	})
	t.Run("reviews.GetAllReviewsPaginated: no movieID or userID specified", func(t *testing.T) {
		_, err := c.GetAllReviewsPaginated(&contracts.GetReviewsRequest{})
		requireBadRequestError(t, err, "either movie_id or user_id must be provided")
	})
	t.Run("reviews.UpdateReview: success", func(t *testing.T) {
		req := &contracts.UpdateReviewRequest{
			ReviewID: review1.ID,
			UserID:   reviewer1.ID,
			Rating:   review1.Rating,
			Title:    review1.Title,
			Content: "I won't call myself a Star Wars fan, but at the same time they were one" +
				" of those who defined my childhood. Having reviewed the cult film by George Lucas," +
				"you seem to be returning to that wonderful time, of course, you look at many things differently. " +
				"But the emotions are still strong, they are different, but very real.",
		}
		err := c.UpdateReview(contracts.NewAuthenticated(req, reviewer1Token))
		require.NoError(t, err)

		review1 = getReview(t, c, review1.ID)
		require.Equal(t, req.Content, review1.Content)
	})
	t.Run("reviews.UpdateReview: mismatch between token and path", func(t *testing.T) {
		req := &contracts.UpdateReviewRequest{
			ReviewID: review2.ID,
			UserID:   reviewer2.ID,
			Rating:   10,
			Content:  ".........",
		}
		err := c.UpdateReview(contracts.NewAuthenticated(req, reviewer1Token))
		requireForbiddenError(t, err, "insufficient permissions")
	})
	t.Run("reviews.DeleteReview: not found", func(t *testing.T) {
		nonExistingID := 10000
		req := &contracts.DeleteReviewRequest{
			ReviewID: nonExistingID,
			UserID:   review3.UserID,
		}
		err := c.DeleteReview(contracts.NewAuthenticated(req, reviewer1Token))
		requireNotFoundError(t, err, "review", "id", nonExistingID)
	})
	t.Run("reviews.DeleteReview: owned by another user", func(t *testing.T) {
		req := &contracts.DeleteReviewRequest{
			ReviewID: review3.ID,
			UserID:   reviewer2.ID,
		}
		err := c.DeleteReview(contracts.NewAuthenticated(req, reviewer2Token))
		requireForbiddenError(t, err, "review with id 3 is not owned by user with id 7")
	})
	t.Run("reviews.DeleteReview: success", func(t *testing.T) {
		req := &contracts.DeleteReviewRequest{
			ReviewID: review1.ID,
			UserID:   reviewer1.ID,
		}
		err := c.DeleteReview(contracts.NewAuthenticated(req, reviewer1Token))
		require.NoError(t, err)

		review1 = getReview(t, c, review1.ID)
		require.Nil(t, review1)
	})
}

func getReview(t *testing.T, c *client.Client, reviewID int) *contracts.Review {
	review, err := c.GetReviewByID(reviewID)
	if err != nil {
		cerr, ok := err.(*client.Error)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cerr.Code)
		return nil
	}
	return review
}
