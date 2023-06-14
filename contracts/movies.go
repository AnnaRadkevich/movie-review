package contracts

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string `json:"description"`
	Version     int    `json:"version"`
}
type GetOrDeleteMovieByIDRequest struct {
	ID int `param:"id" validate:"nonzero"`
}
type GetMoviesRequest struct {
	PaginatedRequest
}
type CreateMovieRequest struct {
	Title       string    `json:"title" validate:"min=1,max=255"`
	ReleaseDate time.Time `json:"release_date" validate:"nonzero"`
	Description string    `json:"description"`
}

type UpdateMovieRequest struct {
	ID          int       `param:"id" validate:"nonzero"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Description string    `json:"description"`
	Version     int       `json:"version" validate:"min=0"`
}
