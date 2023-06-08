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
	FirstName  string     `json:"first_name" validate:"nonzero"`
	MiddleName *string    `json:"middle_name"`
	LastName   string     `json:"last_name" validate:"nonzero"`
	BirthDate  time.Time  `json:"birth_date"`
	BirthPlace *string    `json:"birth_place"`
	DeathDate  *time.Time `json:"death_date"`
	Bio        *string    `json:"bio"`
}

type GetStarByIDRequest struct {
	ID int `param:"starId" validate:"nonzero"`
}
