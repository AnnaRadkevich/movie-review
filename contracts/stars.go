package contracts

import "time"

type Star struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name,omitempty"`
	LastName   string     `json:"last_name"`
	BirthDate  time.Time  `json:"birth_date,omitempty"`
	BirthPlace *string    `json:"birth_place,omitempty"`
	DeathDate  *time.Time `json:"death_date,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type CreateStarRequest struct {
	FirstName  string     `json:"first_name" validate:"min=1,max=50"`
	MiddleName *string    `json:"middle_name,omitempty" validate:"max=50"`
	LastName   string     `json:"last_name" validate:"min=1,max=50"`
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
}
type GetStarsResponse = PaginatedResponse[Star]

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
