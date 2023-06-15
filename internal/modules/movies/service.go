package movies

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

func (s *Service) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	if err := s.repo.CreateMovie(ctx, movie); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"movie created",
		"movie title", movie.Title,
		"movie release date", movie.ReleaseDate)
	return nil
}

func (s *Service) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	return s.repo.GetMovieByID(ctx, id)
}

func (s *Service) GetAllMoviesPaginated(ctx context.Context, offset int, limit int) ([]*MovieDetails, int, error) {
	return s.repo.GetAllMoviesPaginated(ctx, offset, limit)
}

func (s *Service) UpdateMovie(ctx context.Context, movie *MovieDetails) error {
	if err := s.repo.UpdateMovie(ctx, movie); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"movie updated",
		"id", movie.ID)
	return nil
}

func (s *Service) DeleteMovie(ctx context.Context, id int) error {
	if err := s.repo.DeleteMovie(ctx, id); err != nil {
		return err
	}
	log.FromContext(ctx).Info(
		"movie deleted",
		"id", id)
	return nil
}
