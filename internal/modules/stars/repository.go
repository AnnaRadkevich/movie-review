package stars

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateStar(ctx context.Context, star *Star) error {
	queryString := `INSERT INTO stars 
(first_name, middle_name, last_name, birth_date, birth_place, death_date, bio)
 VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING
	id, created_at, deleted_at`
	row := r.db.QueryRow(ctx, queryString, star.FirstName, star.MiddleName, star.LastName, star.BirthDate,
		star.BirthPlace, star.DeathDate, star.Bio)
	err := row.Scan(&star.ID, &star.CreatedAt, &star.DeletedAt)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetStarByID(ctx context.Context, id int) (*Star, error) {
	var star Star
	queryString := `
	SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at
	FROM stars
	WHERE id = $1 and deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(&star.ID,
		&star.FirstName,
		&star.MiddleName,
		&star.LastName,
		&star.BirthDate,
		&star.BirthPlace,
		&star.DeathDate,
		&star.Bio,
		&star.CreatedAt,
		&star.DeletedAt)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("star", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &star, nil
}
