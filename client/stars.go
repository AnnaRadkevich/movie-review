package client

import (
	"github.com/cloudmachinery/movie-reviews/contracts"
)

func (c *Client) CreateStar(req *contracts.AuthenticadedRequest[*contracts.CreateStarRequest]) (*contracts.StarDetails, error) {
	var star contracts.StarDetails
	_, err := c.client.R().SetResult(&star).SetAuthToken(req.AccessToken).
		SetBody(req.Request).Post(c.path("/api/stars"))
	return &star, err
}

func (c *Client) GetStarByID(id int) (*contracts.StarDetails, error) {
	var star contracts.StarDetails

	_, err := c.client.R().SetResult(&star).Get(c.path("/api/stars/%d", id))
	return &star, err
}

func (c *Client) GetStars(req *contracts.GetStarsRequest) (*contracts.PaginatedResponse[contracts.Star], error) {
	var stars contracts.PaginatedResponse[contracts.Star]
	_, err := c.client.R().SetResult(&stars).SetQueryParams(req.ToQueryParams()).Get(c.path("/api/stars"))
	return &stars, err
}

func (c *Client) UpdateStar(req *contracts.AuthenticadedRequest[*contracts.UpdateStarRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Put(c.path("/api/stars/%d", req.Request.ID))
	return err
}

func (c *Client) DeleteStar(req *contracts.AuthenticadedRequest[*contracts.GetOrDeleteStarByIDRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Delete(c.path("/api/stars/%d", req.Request.ID))
	return err
}
