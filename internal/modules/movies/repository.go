package movies

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

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	queryString := `INSERT INTO movies 
(title,release_date,description) 
VALUES ($1,$2,$3)
	RETURNING id,created_at,deleted_at`
	row := r.db.QueryRow(ctx, queryString, movie.Title, movie.ReleaseDate, movie.Description)
	err := row.Scan(&movie.ID, &movie.CreatedAt, &movie.DeletedAt)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	var movie MovieDetails
	queryString := `SELECT id,title,release_date,created_at,deleted_at,description,version
 FROM movies WHERE id=$1 AND deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate,
		&movie.CreatedAt, &movie.DeletedAt, &movie.Description, &movie.Version)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("movie", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &movie, nil
}

func (r *Repository) GetAllMoviesPaginated(ctx context.Context, offset int, limit int) ([]*MovieDetails, int, error) {
	b := &pgx.Batch{}
	b.Queue(`SELECT id,title,release_date,created_at,deleted_at FROM movies
		WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	b.Queue(`SELECT COUNT(*) FROM movies WHERE deleted_at IS NULL`)
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var movies []*MovieDetails
	for rows.Next() {
		var movie MovieDetails
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.CreatedAt,
			&movie.DeletedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return movies, total, err
}

func (r *Repository) UpdateMovie(ctx context.Context, movie *MovieDetails) error {
	n, err := r.db.Exec(ctx, `UPDATE movies 
	SET 
		title = $1,
		release_date = $2,
		description = $3,
		version = version + 1
	WHERE id = $4 and deleted_at IS NULL and version = $5`,
		movie.Title, movie.ReleaseDate, movie.Description, movie.ID, movie.Version)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		_, err := r.GetMovieByID(ctx, movie.ID)
		if err != nil {
			return err
		}
		return apperrors.VersionMismatch("movie", "id", movie.ID, movie.Version)
	}

	return nil
}

func (r *Repository) DeleteMovie(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `UPDATE movies SET deleted_at = NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", id)
	}
	return nil
}
