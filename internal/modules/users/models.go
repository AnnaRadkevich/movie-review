package users

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Bio       *string    `json:"bio,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

type UserWithPassword struct {
	*User
	PasswordHash string
}
