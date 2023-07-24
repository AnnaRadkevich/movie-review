package client

import "github.com/RadkevichAnn/movie-reviews/contracts"

func (c *Client) CreateReview(req *contracts.AuthenticadedRequest[*contracts.CreateReviewRequest]) (*contracts.Review, error) {
	var review contracts.Review
	_, err := c.client.R().SetResult(&review).SetAuthToken(req.AccessToken).
		SetBody(req.Request).Post(c.path("/api/users/%d/reviews", req.Request.UserID))
	return &review, err
}

func (c *Client) GetReviewByID(id int) (*contracts.Review, error) {
	var review contracts.Review
	_, err := c.client.R().SetResult(&review).Get(c.path("/api/reviews/%d", id))
	return &review, err
}

func (c *Client) GetAllReviewsPaginated(req *contracts.GetReviewsRequest) (*contracts.PaginatedResponse[contracts.Review], error) {
	var reviews contracts.PaginatedResponse[contracts.Review]

	_, err := c.client.R().SetResult(&reviews).SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/reviews"))
	return &reviews, err
}

func (c *Client) UpdateReview(req *contracts.AuthenticadedRequest[*contracts.UpdateReviewRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).Put(c.path("/api/users/%d/reviews/%d", req.Request.UserID, req.Request.ReviewID))
	return err
}

func (c *Client) DeleteReview(req *contracts.AuthenticadedRequest[*contracts.DeleteReviewRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).SetBody(req.Request).
		Delete(c.path("/api/users/%d/reviews/%d", req.Request.UserID, req.Request.ReviewID))
	return err
}
