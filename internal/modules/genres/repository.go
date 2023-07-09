package genres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllGenres(ctx context.Context) ([]*Genre, error) {
	queryString := `SELECT id, name FROM genres `
	rows, err := r.db.Query(ctx, queryString)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()
	return pgx.CollectRows[*Genre](rows, pgx.RowToAddrOfStructByPos[Genre])
}

func (r *Repository) GetGenreById(ctx context.Context, id int) (*Genre, error) {
	queryString := `SELECT id, name FROM genres WHERE id = $1`
	row := r.db.QueryRow(ctx, queryString, id)

	var genre Genre
	err := row.Scan(&genre.ID, &genre.Name)
	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("genre", "id", id)
	case err != nil:
		return nil, apperrors.Internal(err)
	}
	return &genre, nil
}

func (r *Repository) CreateGenre(ctx context.Context, name string) (*Genre, error) {
	queryString := "INSERT INTO genres (name) VALUES ($1) returning id, name;"
	row := r.db.QueryRow(ctx, queryString, name)

	var genre Genre
	err := row.Scan(&genre.ID, &genre.Name)
	if dbx.IsUniqueViolation(err, "name") {
		return nil, apperrors.AlreadyExists("genre", "name", name)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &genre, nil
}

func (r *Repository) UpdateGenre(ctx context.Context, id int, name string) error {
	n, err := r.db.Exec(ctx, "UPDATE genres SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}
	return nil
}

func (r *Repository) DeleteGenre(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, "DELETE FROM genres WHERE id = $1", id)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}
	return nil
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, id int) ([]*MovieGenreRelation, error) {
	queryString := `SELECT movie_id,genre_id,order_no FROM movie_genres WHERE movie_id = $1`
	rows, err := r.db.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var relations []*MovieGenreRelation
	for rows.Next() {
		var relation MovieGenreRelation
		err := rows.Scan(
			&relation.MovieID,
			&relation.GenreID,
			&relation.OrderNo)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		relations = append(relations, &relation)
	}
	return relations, nil
}

func (r *Repository) GetGenresByMovieID(ctx context.Context, id int) ([]*Genre, error) {
	queryString := `SELECT g.id, g.name
	FROM genres g
	INNER JOIN movie_genres mg on mg.genre_id = g.id	
	WHERE mg.movie_id = $1
	ORDER BY mg.order_no`
	rows, err := r.db.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return pgx.CollectRows[*Genre](rows, pgx.RowToAddrOfStructByPos[Genre])
}
