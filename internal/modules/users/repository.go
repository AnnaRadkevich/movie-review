package users

import (
	"context"

	"github.com/RadkevichAnn/movie-reviews/internal/apperrors"
	"github.com/RadkevichAnn/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *UserWithPassword) error {
	queryString := "INSERT INTO users (username, email, pass_hash, role) VALUES ($1, $2, $3, $4) returning id, created_at"
	err := r.db.QueryRow(ctx, queryString, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return apperrors.AlreadyExists("user", "email", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return apperrors.AlreadyExists("user", "username", user.Username)
	case err != nil:
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetExistingUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	queryString := `
	SELECT id, username, email, pass_hash, role, created_at, deleted_at, bio
	FROM users
	WHERE email = $1 AND deleted_at IS NULL;`

	user := newUserWithPassword()

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
	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, userId int) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", userId)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", userId)
	}
	return nil
}

func (r *Repository) GetUserById(ctx context.Context, userId int) (*User, error) {
	var user User
	query := "SELECT id, username, email,  role, bio FROM users WHERE id = $1 AND deleted_at IS NULL  "
	row := r.db.QueryRow(ctx, query, userId)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Bio)
	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", userId)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &user, nil
}

func (r *Repository) UpdateBio(ctx context.Context, userId int, bio string) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET bio = $1 WHERE id = $2 AND deleted_at IS NULL", bio, userId)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", userId)
	}
	return nil
}

func (r *Repository) UpdateRole(ctx context.Context, userId int, role string) error {
	n, err := r.db.Exec(ctx, "UPDATE users SET role = $1 WHERE id = $2 AND deleted_at IS NULL", role, userId)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", userId)
	}
	return nil
}

func (r *Repository) GetExistingUserByUsername(ctx context.Context, username string) (*User, error) {
	queryString := `
	SELECT id, username, email, role, created_at, deleted_at, bio
	FROM users
	WHERE username = $1 and deleted_at IS NULL;`

	user := User{}

	row := r.db.QueryRow(ctx, queryString, username)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("user", "username", username)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &user, nil
}
