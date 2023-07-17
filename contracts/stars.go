package contracts

import (
	"strconv"
	"time"
)

type Star struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate time.Time  `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
type StarDetails struct {
	Star
	MiddleName *string `json:"middle_name,omitempty"`
	BirthPlace *string `json:"birth_place,omitempty"`
	Bio        *string `json:"bio,omitempty"`
}

type CreateStarRequest struct {
	FirstName  string     `json:"first_name" validate:"min=1,max=50"`
	MiddleName *string    `json:"middle_name,omitempty" validate:"max=50"`
	LastName   string     `json:"last_name" validate:"max=50"`
	BirthDate  time.Time  `json:"birth_date" validate:"nonzero"`
	BirthPlace *string    `json:"birth_place,omitempty" validate:"max=100"`
	DeathDate  *time.Time `json:"death_date,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
}

type GetOrDeleteStarByIDRequest struct {
	ID int `param:"id" validate:"nonzero"`
}

type GetStarsRequest struct {
	PaginatedRequest
	MovieID *int `query:"movieId"`
}

func (r *GetStarsRequest) ToQueryParams() map[string]string {
	params := r.PaginatedRequest.ToQueryParams()
	if r.MovieID != nil {
		params["movieId"] = strconv.Itoa(*r.MovieID)
	}
	return params
}

type GetStarsResponse = PaginatedResponse[StarDetails]

type UpdateStarRequest struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name,omitempty"`
	LastName   string     `json:"last_name"`
	BirthDate  time.Time  `json:"birth_date,omitempty"`
	BirthPlace *string    `json:"birth_place,omitempty"`
	DeathDate  *time.Time `json:"death_date,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
}
