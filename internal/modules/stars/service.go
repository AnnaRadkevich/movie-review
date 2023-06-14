package stars

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

func (s *Service) CreateStar(ctx context.Context, star *StarDetails) error {
	if err := s.repo.CreateStar(ctx, star); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"star created",
		"star first name", star.FirstName,
		"star last name", star.LastName)
	return nil
}

func (s *Service) GetStarByID(ctx context.Context, id int) (*StarDetails, error) {
	return s.repo.GetStarByID(ctx, id)
}

func (s *Service) GetAllStarsPaginated(ctx context.Context, offset int, limit int) ([]*StarDetails, int, error) {
	return s.repo.GetAllStarsPaginated(ctx, offset, limit)
}

func (s *Service) UpdateStar(ctx context.Context, star *StarDetails) error {
	if err := s.repo.UpdateStar(ctx, star); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"star updated",
		"id", star.ID)
	return nil
}

func (s *Service) DeleteStar(ctx context.Context, id int) error {
	if err := s.repo.DeleteStar(ctx, id); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"star deleted",
		"id", id)
	return nil
}
