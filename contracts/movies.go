package contracts

import (
	"strconv"
	"time"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string         `json:"description"`
	Version     int            `json:"version"`
	Genres      []*Genre       `json:"genres"`
	Cast        []*MovieCredit `json:"cast"`
}
type MovieCredit struct {
	Star    Star    `json:"star"`
	Role    string  `json:"role"`
	Details *string `json:"details,omitempty"`
}
type MovieCreditInfo struct {
	StarID  int     `json:"star_id"`
	Role    string  `json:"role"`
	Details *string `json:"details"`
}
type GetOrDeleteMovieByIDRequest struct {
	ID int `param:"id" validate:"nonzero"`
}
type GetMoviesRequest struct {
	PaginatedRequest
	StarID     *int    `query:"starId"`
	SearchTerm *string `query:"q"`
}

func (r *GetMoviesRequest) ToQueryParams() map[string]string {
	params := r.PaginatedRequest.ToQueryParams()
	if r.StarID != nil {
		params["starId"] = strconv.Itoa(*r.StarID)
	}
	if r.SearchTerm != nil {
		params["q"] = *r.SearchTerm
	}
	return params
}

type CreateMovieRequest struct {
	Title       string             `json:"title" validate:"min=1,max=255"`
	ReleaseDate time.Time          `json:"release_date" validate:"nonzero"`
	Description string             `json:"description"`
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}

type UpdateMovieRequest struct {
	ID          int                `param:"id" validate:"nonzero"`
	Title       string             `json:"title"`
	ReleaseDate time.Time          `json:"release_date"`
	Description string             `json:"description"`
	Version     int                `json:"version" validate:"min=0"`
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}
