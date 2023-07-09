package reviews

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateReview(ctx context.Context, review *Review) error {
	if err := s.repo.CreateReview(ctx, review); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"review created")
	return nil
}

func (s *Service) GetReviewByID(ctx context.Context, id int) (*Review, error) {
	r, err := s.repo.GetReviewByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (s *Service) GetAllReviewsPaginated(ctx context.Context, movieID, userID *int, offset int, limit int) ([]*Review, int, error) {
	return s.repo.GetAllReviewsPaginated(ctx, movieID, userID, offset, limit)
}

func (s *Service) UpdateReview(ctx context.Context, reviewID, userID int, title, content string, rating int) error {
	if err := s.repo.UpdateReview(ctx, reviewID, userID, title, content, rating); err != nil {
		return err
	}
	log.FromContext(ctx).Info("review updated",
		"id", reviewID)
	return nil
}

func (s *Service) DeleteReview(ctx context.Context, reviewID, userID int) error {
	if err := s.repo.DeleteReview(ctx, reviewID, userID); err != nil {
		return err
	}
	log.FromContext(ctx).Info("review deleted",
		"id", reviewID)
	return nil
}
