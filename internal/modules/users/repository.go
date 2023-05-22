package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *UserWithPassword) error {
	queryString := "INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) returning id, created_at, role"
	err := r.db.QueryRow(ctx, queryString, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.Role)

	return err
}

func (r *Repository) GetExistingUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	queryString := `
	SELECT id, username, email, pass_hash, role, created_at, deleted_at, bio
	FROM users
	WHERE email = $1 AND deleted_at IS NULL;`

	user := UserWithPassword{
		User: &User{},
	}

	row := r.db.QueryRow(ctx, queryString, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if err != nil {
		return nil, fmt.Errorf("scan repo: %w", err)
	}

	return &user, nil
}
