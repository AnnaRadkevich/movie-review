package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) GetAllGenres() ([]*contracts.Genre, error) {
	var genres []*contracts.Genre

	_, err := c.client.R().SetResult(&genres).Get(c.path("/api/genres"))
	return genres, err
}

func (c *Client) GetGenreById(id int) (*contracts.Genre, error) {
	var g contracts.Genre

	_, err := c.client.R().SetResult(&g).Get(c.path("/api/genres/%d", id))
	return &g, err
}

func (c *Client) CreateGenre(req *contracts.AuthenticadedRequest[*contracts.CreateGenreRequest]) (*contracts.Genre, error) {
	var genre contracts.Genre
	_, err := c.client.R().
		SetResult(&genre).
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Post(c.path("/api/genres"))

	return &genre, err
}

func (c *Client) UpdateGenre(req *contracts.AuthenticadedRequest[*contracts.UpdateGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Put(c.path("/api/genres/%d", req.Request.GenreId))

	return err
}

func (c *Client) DeleteGenre(req *contracts.AuthenticadedRequest[*contracts.DeleteGenreRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Delete(c.path("/api/genres/%d", req.Request.GenreId))
	return err
}
