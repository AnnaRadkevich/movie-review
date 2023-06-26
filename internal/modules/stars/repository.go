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

func (r *Repository) GetAllStarsPaginated(ctx context.Context, movieID *int, offset int, limit int) ([]*StarDetails, int, error) {
	b := &pgx.Batch{}

	selectQuery := dbx.StatementBuilder.
		Select("id, first_name, last_name, birth_date, death_date, created_at,deleted_at").
		From("stars").
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	countQuery := dbx.StatementBuilder.
		Select("count(*)").
		From("stars").
		Where("deleted_at IS NULL")

	if movieID != nil {
		selectQuery = selectQuery.
			Join("movie_stars on stars.id = movie_stars.star_id").
			Where("movie_stars.movie_id = ?", movieID)

		countQuery = countQuery.
			Join("movie_stars on stars.id = movie_stars.star_id").
			Where("movie_stars.movie_id = ?", movieID)
	}
	if err := dbx.QueueBatchSelect(b, selectQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	if err := dbx.QueueBatchSelect(b, countQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

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

func (r *Repository) GetCastByMovieID(ctx context.Context, movieID int) ([]*MovieCredit, error) {
	queryString := `SELECT s.id, s.first_name, s.last_name, s.birth_date, s.death_date, s.created_at, ms.role, ms.details 
			FROM stars s
			INNER JOIN movie_stars ms ON ms.star_id = s.id
			WHERE ms.movie_id = $1
			ORDER BY ms.order_no`
	rows, err := r.db.Query(ctx, queryString, movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	var cast []*MovieCredit
	for rows.Next() {
		var mc MovieCredit
		err = rows.Scan(
			&mc.Star.ID,
			&mc.Star.FirstName,
			&mc.Star.LastName,
			&mc.Star.BirthDate,
			&mc.Star.DeathDate,
			&mc.Star.CreatedAt,
			&mc.Role,
			&mc.Details,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		cast = append(cast, &mc)
	}
	return cast, nil
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, id int) ([]*MovieStarRelation, error) {
	queryString := `
	SELECT movie_id, star_id, role, details, order_no
	FROM movie_stars
	WHERE movie_id = $1`
	q := dbx.FromContext(ctx, r.db)
	rows, err := q.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	var relations []*MovieStarRelation
	for rows.Next() {
		var r MovieStarRelation
		err = rows.Scan(
			&r.MovieID,
			&r.StarID,
			&r.Role,
			&r.Details,
			&r.OrderNo,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		relations = append(relations, &r)
	}
	return relations, nil
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
