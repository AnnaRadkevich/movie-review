package reviews

import (
	"context"
	"fmt"

	"github.com/RadkevichAnn/movie-reviews/internal/apperrors"
	"github.com/RadkevichAnn/movie-reviews/internal/dbx"
	"github.com/RadkevichAnn/movie-reviews/internal/modules/movies"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db         *pgxpool.Pool
	moviesRepo *movies.Repository
}

func NewRepository(db *pgxpool.Pool, moviesRepo *movies.Repository) *Repository {
	return &Repository{
		db:         db,
		moviesRepo: moviesRepo,
	}
}

func (r *Repository) CreateReview(ctx context.Context, review *Review) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		if err := r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}

		// Insert Review
		err := tx.QueryRow(ctx, `INSERT INTO reviews (movie_id, user_id, title, content, rating) VALUES ($1, $2, $3, $4, $5) returning id, created_at`,
			review.MovieID, review.UserID, review.Title, review.Content, review.Rating).
			Scan(&review.ID, &review.CreatedAt)
		switch {
		case dbx.IsUniqueViolation(err, ""):
			return apperrors.AlreadyExists("review", "(movie_id,user_id)", fmt.Sprintf("(%d,%d)", review.MovieID, review.UserID))
		case err != nil:
			return apperrors.Internal(err)
		}
		return r.recalculateMovieRating(ctx, review.MovieID)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}
	return nil
}

func (r *Repository) GetReviewByID(ctx context.Context, reviewID int) (*Review, error) {
	var review Review

	err := r.db.QueryRow(ctx, `SELECT id, movie_id, user_id, title, content, rating, created_at FROM reviews where deleted_at is null and id = $1`,
		reviewID).
		Scan(&review.ID, &review.MovieID, &review.UserID, &review.Title, &review.Content, &review.Rating, &review.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("review", "id", reviewID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &review, nil
}

func (r *Repository) GetAllReviewsPaginated(ctx context.Context, movieID *int, userID *int, offset int, limit int) ([]*Review, int, error) {
	selectQuery := dbx.StatementBuilder.
		Select("id", "movie_id", "user_id", "title", "content", "rating", "created_at").
		From("reviews").
		Where("deleted_at is NULL").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	queryTotal := dbx.StatementBuilder.
		Select("COUNT(*)").
		From("reviews").
		Where("deleted_at is NULL")
	if movieID != nil {
		selectQuery = selectQuery.Where("movie_id = ?", *movieID)
		queryTotal = queryTotal.Where("movie_id = ?", *movieID)
	}
	if userID != nil {
		selectQuery = selectQuery.Where("user_id = ?", *userID)
		queryTotal = queryTotal.Where("user_id = ?", *userID)
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

	var reviews []*Review
	for rows.Next() {
		var review Review
		if err := rows.Scan(
			&review.ID,
			&review.MovieID,
			&review.UserID,
			&review.Title,
			&review.Content,
			&review.Rating,
			&review.CreatedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		reviews = append(reviews, &review)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return reviews, total, err
}

func (r *Repository) UpdateReview(ctx context.Context, reviewID, userID int, title, content string, rating int) error {
	review, err := r.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}
	err = dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		if err = r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}
		var n pgconn.CommandTag
		n, err = r.db.Exec(ctx, "UPDATE reviews SET title = $1, content = $2, rating = $3 WHERE deleted_at IS NULL AND id = $4 AND user_id = $5",
			title, content, rating, reviewID, userID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if n.RowsAffected() == 0 {
			return r.specifyModificationError(ctx, reviewID, userID)
		}

		return r.recalculateMovieRating(ctx, review.MovieID)
	})
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) DeleteReview(ctx context.Context, reviewID, userID int) error {
	review, err := r.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}
	err = dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		if err = r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}
		var n pgconn.CommandTag
		n, err = r.db.Exec(ctx,
			`UPDATE reviews SET deleted_at = now() WHERE deleted_at IS NULL AND id = $1 AND user_id = $2`,
			reviewID, userID)

		if n.RowsAffected() == 0 {
			return r.specifyModificationError(ctx, reviewID, userID)
		}

		return r.recalculateMovieRating(ctx, review.MovieID)
	})

	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) specifyModificationError(ctx context.Context, reviewID, userID int) error {
	// Review is not found by reviewID and userID then there are two possibilities:
	// 1. Review with reviewID does not exist
	// 2. Review with reviewID exists, but it is not owned by userID
	review, err := r.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return apperrors.Forbidden(fmt.Sprintf("review with id %d is not owned by user with id %d", reviewID, userID))
	}

	// If we got here, then something is wrong
	return apperrors.Internal(fmt.Errorf("unexpected error creating/updating review with id %d", reviewID))
}

func (r *Repository) recalculateMovieRating(ctx context.Context, movieID int) error {
	q := dbx.FromContext(ctx, r.db)
	n, err := q.Exec(ctx, `UPDATE movies SET avg_rating = (SELECT AVG(rating) FROM reviews WHERE deleted_at IS NULL and movie_id = $1) where id = $1`, movieID)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", movieID)
	}
	return nil
}
