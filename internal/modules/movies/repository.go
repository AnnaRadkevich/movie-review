package movies

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/modules/stars"

	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"
	"github.com/cloudmachinery/movie-reviews/internal/slices"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db               *pgxpool.Pool
	genresRepository *genres.Repository
	starsRepository  *stars.Repository
}

func NewRepository(db *pgxpool.Pool, genresRepository *genres.Repository, starsRepository *stars.Repository) *Repository {
	return &Repository{
		db:               db,
		genresRepository: genresRepository,
		starsRepository:  starsRepository,
	}
}

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		// Insert movies
		queryString := `INSERT INTO movies 
						(title,release_date,description) 
						VALUES ($1,$2,$3)
							RETURNING id,created_at,deleted_at`
		err := tx.QueryRow(ctx, queryString, movie.Title, movie.ReleaseDate, movie.Description).
			Scan(&movie.ID, &movie.CreatedAt, &movie.DeletedAt)
		if err != nil {
			return err
		}

		// Insert genres
		nextGenres := slices.MapIndex(movie.Genres, func(i int, g *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: movie.ID,
				GenreID: g.ID,
				OrderNo: i,
			}
		})
		if err = r.UpdateGenres(ctx, []*genres.MovieGenreRelation{}, nextGenres); err != nil {
			return err
		}

		// Insert cast
		nextCast := slices.MapIndex(movie.Cast, func(i int, c *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: movie.ID,
				StarID:  c.Star.ID,
				Role:    c.Role,
				Details: c.Details,
				OrderNo: i,
			}
		})
		return r.UpdateCast(ctx, nil, nextCast)
	})
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	var movie MovieDetails
	queryString := `SELECT id,title,release_date,avg_rating,created_at,deleted_at,description,version
 FROM movies WHERE id=$1 AND deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.AvgRating,
		&movie.CreatedAt, &movie.DeletedAt, &movie.Description, &movie.Version)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("movie", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &movie, nil
}

func (r *Repository) GetAllMoviesPaginated(ctx context.Context, searchTerm *string, sortByRating *string, starID *int, offset int, limit int) ([]*MovieDetails, int, error) {
	selectQuery := dbx.StatementBuilder.
		Select("id,title,release_date,avg_rating,created_at,deleted_at").
		From("movies").
		Where("deleted_at IS NULL").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	queryTotal := dbx.StatementBuilder.
		Select("COUNT(*)").
		From("movies").
		Where("deleted_at IS NULL")

	if sortByRating != nil {
		selectQuery = selectQuery.
			OrderByClause("avg_rating" + *sortByRating)
	}
	if starID != nil {
		selectQuery = selectQuery.
			Join("movie_stars on movies.id = movie_stars.movie_id").
			Where("star_id = ?", starID)

		queryTotal = queryTotal.
			Join("movie_stars on movies.id = movie_stars.movie_id").
			Where("star_id = ?", starID)
	}
	if searchTerm != nil {
		selectQuery = selectQuery.
			Where("search_vector @@ to_tsquery('english', ?)", *searchTerm).
			OrderByClause("ts_rank_cd(search_vector, to_tsquery('english', ?)) DESC", *searchTerm)

		queryTotal = queryTotal.
			Where("search_vector @@ to_tsquery('english', ?)", *searchTerm)
	}

	b := &pgx.Batch{}
	if err := dbx.QueueBatchSelect(b, selectQuery); err != nil {
		return nil, 0, err
	}
	if err := dbx.QueueBatchSelect(b, queryTotal); err != nil {
		return nil, 0, err
	}
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
			&movie.AvgRating,
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
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		n, err := tx.
			Exec(
				ctx,
				`UPDATE movies 
			SET version = version + 1, 
			title = $1,
			description = $2, 
			release_date = $3 
			WHERE id = $4 
			AND version = $5;`,
				movie.Title,
				movie.Description,
				movie.ReleaseDate,
				movie.ID,
				movie.Version,
			)
		if err != nil {
			return apperrors.Internal(err)
		}

		if n.RowsAffected() == 0 {
			_, err = r.GetMovieByID(ctx, movie.ID)
			if err != nil {
				return err
			}

			return apperrors.VersionMismatch("movie", "id", movie.ID, movie.Version)
		}

		currentGenres, err := r.genresRepository.GetRelationsByMovieID(ctx, movie.ID)
		if err != nil {
			return err
		}

		nextGenres := slices.MapIndex(movie.Genres, func(i int, g *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				GenreID: g.ID,
				MovieID: movie.ID,
				OrderNo: i,
			}
		})

		err = r.UpdateGenres(ctx, currentGenres, nextGenres)
		if err != nil {
			return err
		}
		currentCast, err := r.starsRepository.GetRelationsByMovieID(ctx, movie.ID)
		if err != nil {
			return err
		}
		nextCast := slices.MapIndex(movie.Cast, func(i int, c *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: movie.ID,
				StarID:  c.Star.ID,
				Role:    c.Role,
				Details: c.Details,
				OrderNo: i,
			}
		})
		if err != nil {
			return err
		}

		return r.UpdateCast(ctx, currentCast, nextCast)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) DeleteMovie(ctx context.Context, id int) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		n, err := r.db.Exec(ctx, `UPDATE movies SET deleted_at = NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
		if err != nil {
			return apperrors.Internal(err)
		}
		if n.RowsAffected() == 0 {
			return apperrors.NotFound("movie", "id", id)
		}

		current, err := r.genresRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}
		err = r.UpdateGenres(ctx, current, []*genres.MovieGenreRelation{})
		if err != nil {
			return err
		}
		currentCast, err := r.starsRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}
		return r.UpdateCast(ctx, currentCast, []*stars.MovieStarRelation{})
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) Lock(ctx context.Context, tx pgx.Tx, movieID int) error {
	n, err := r.db.Exec(ctx, `SELECT 1 FROM movies WHERE deleted_at IS NULL AND id = $1 for update`, movieID)
	if err != nil {
		return err
	}
	if n.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", movieID)
	}
	return nil
}

func (r *Repository) UpdateGenres(ctx context.Context, current, next []*genres.MovieGenreRelation) error {
	q := dbx.FromContext(ctx, r.db)
	addFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx, `INSERT INTO movie_genres (movie_id, genre_id, order_no) VALUES ($1, $2, $3)`,
			mgo.MovieID, mgo.GenreID, mgo.OrderNo)
		return err
	}
	removeFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx, `DELETE FROM movie_genres WHERE movie_id = $1 and genre_id = $2`,
			mgo.MovieID, mgo.GenreID)
		return err
	}
	return dbx.AdjustRelations(current, next, addFunc, removeFunc)
}

func (r *Repository) UpdateCast(ctx context.Context, current, next []*stars.MovieStarRelation) error {
	q := dbx.FromContext(ctx, r.db)
	addFunc := func(msr *stars.MovieStarRelation) error {
		_, err := q.Exec(ctx, `INSERT INTO movie_stars (movie_id, star_id, role, details, order_no) 
				VALUES ($1, $2, $3, $4, $5)`,
			msr.MovieID,
			msr.StarID,
			msr.Role,
			msr.Details,
			msr.OrderNo,
		)
		return err
	}
	removeFunc := func(msr *stars.MovieStarRelation) error {
		_, err := q.Exec(ctx, `DELETE FROM movie_stars 
				WHERE movie_id = $1 
				AND star_id = $2 
				AND role = $3`,
			msr.MovieID,
			msr.StarID,
			msr.Role,
		)
		return err
	}
	return dbx.AdjustRelations(current, next, addFunc, removeFunc)
}
