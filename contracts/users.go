package contracts

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Bio       *string    `json:"bio,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DeleteOrGetRequest struct {
	UserId int `param:"userId"`
}

type UpdateRequest struct {
	UserId int     `param:"userId"`
	Bio    *string `json:"bio"`
}

type UpdateRoleRequest struct {
	UserID int    `param:"userId" validate:"nonzero"`
	Role   string `param:"role" validate:"role"`
}

type GetUserByUsernameRequest struct {
	Username string `param:"username" validate:"nonzero"`
}
