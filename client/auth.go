package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) RegisterUser(req *contracts.RegisterUserRequest) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().SetBody(req).SetResult(&u).Post(c.path("/api/auth/register"))

	return &u, err
}

func (c *Client) LoginUser(req *contracts.LoginUserRequest) (*contracts.LoginUserResponse, error) {
	var res contracts.LoginUserResponse

	_, err := c.client.R().SetBody(req).SetResult(&res).Post(c.path("/api/auth/login"))

	return &res, err
}