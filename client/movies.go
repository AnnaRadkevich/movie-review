package client

import (
	"github.com/cloudmachinery/movie-reviews/contracts"
)

func (c *Client) CreateMovie(req *contracts.AuthenticadedRequest[*contracts.CreateMovieRequest]) (*contracts.MovieDetails, error) {
	var movie contracts.MovieDetails
	_, err := c.client.R().SetResult(&movie).SetAuthToken(req.AccessToken).
		SetBody(req.Request).Post(c.path("/api/movies"))
	return &movie, err
}

func (c *Client) GetMovieByID(id int) (*contracts.MovieDetails, error) {
	var movie contracts.MovieDetails
	_, err := c.client.R().SetResult(&movie).Get(c.path("/api/movies/%d", id))
	return &movie, err
}

func (c *Client) GetMovies(req *contracts.GetMoviesRequest) (*contracts.PaginatedResponse[contracts.Movie], error) {
	var movies contracts.PaginatedResponse[contracts.Movie]

	_, err := c.client.R().SetResult(&movies).SetQueryParams(req.ToQueryParams()).Get(c.path("/api/movies"))
	return &movies, err
}

func (c *Client) UpdateMovie(req *contracts.AuthenticadedRequest[*contracts.UpdateMovieRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Put(c.path("/api/movies/%d", req.Request.ID))
	return err
}

func (c *Client) DeleteMovie(req *contracts.AuthenticadedRequest[*contracts.GetOrDeleteMovieByIDRequest]) error {
	_, err := c.client.R().SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Delete(c.path("/api/movies/%d", req.Request.ID))
	return err
}
