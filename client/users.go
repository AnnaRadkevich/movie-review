package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) GetUser(userId int) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().SetResult(&u).Get(c.path("/api/users/%d", userId))

	return &u, err
}

func (c *Client) GetUserByUserName(Username string) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().SetResult(&u).Get(c.path("/api/users/username/%s", Username))

	return &u, err
}

func (c *Client) UpdateUser(req *contracts.AuthenticadedRequest[*contracts.UpdateRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/users/%d", req.Request.UserId))
	return err
}

func (c *Client) Delete(req *contracts.AuthenticadedRequest[*contracts.DeleteOrGetRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Delete(c.path("/api/users/%d", req.Request.UserId))
	return err
}

func (c *Client) UpdateUserRole(req *contracts.AuthenticadedRequest[*contracts.UpdateRoleRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).Put(c.path("/api/users/%d/role/%s", req.Request.UserID, req.Request.Role))
	return err
}
