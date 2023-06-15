package stars

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) CreateStar(ctx context.Context, star *StarDetails) error {
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

func (r *Repository) GetStarByID(ctx context.Context, id int) (*StarDetails, error) {
	var star StarDetails
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

func (r *Repository) GetAllStarsPaginated(ctx context.Context, offset int, limit int) ([]*StarDetails, int, error) {
	b := &pgx.Batch{}
	b.Queue("SELECT id, first_name, last_name, birth_date, death_date, created_at, deleted_at FROM stars WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	b.Queue("SELECT COUNT(*) FROM stars WHERE deleted_at IS NULL")
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var stars []*StarDetails
	for rows.Next() {
		var star StarDetails
		if err := rows.
			Scan(
				&star.ID,
				&star.FirstName,
				&star.LastName,
				&star.BirthDate,
				&star.DeathDate,
				&star.CreatedAt,
				&star.DeletedAt,
			); err != nil {
			return nil, 0, apperrors.Internal(err)
		}

		stars = append(stars, &star)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return stars, total, err
}

func (r *Repository) UpdateStar(ctx context.Context, star *StarDetails) error {
	n, err := r.db.Exec(ctx, `UPDATE stars 
			SET first_name = $1, 
			middle_name = $2, 
			last_name = $3, 
			birth_date = $4, 
			birth_place = $5, 
			death_date = $6, 
			bio = $7 
			WHERE id = $8`,
		star.FirstName,
		star.MiddleName,
		star.LastName,
		star.BirthDate,
		star.BirthPlace,
		star.DeathDate,
		star.Bio,
		star.ID)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("star", "id", star.ID)
	}

	return nil
}

func (r *Repository) DeleteStar(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `UPDATE stars SET deleted_at  = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("star", "id", id)
	}
	return nil
}
