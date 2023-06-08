package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) CreateStar(req *contracts.AuthenticadedRequest[*contracts.CreateStarRequest]) (*contracts.Star, error) {
	var star contracts.Star
	_, err := c.client.R().SetResult(&star).SetAuthToken(req.AccessToken).
		SetBody(req.Request).Post(c.path("/api/stars"))
	return &star, err
}

func (c *Client) GetStarByID(starId int) (*contracts.Star, error) {
	var star contracts.Star

	_, err := c.client.R().SetResult(&star).Get(c.path("/api/stars/%d", starId))
	return &star, err
}
