package movies

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/cloudmachinery/movie-reviews/internal/modules/stars"

	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"

	"github.com/cloudmachinery/movie-reviews/internal/log"
)

type Service struct {
	repo          *Repository
	genresService *genres.Service
	starService   *stars.Service
}

func NewService(repo *Repository, genresService *genres.Service, starService *stars.Service) *Service {
	return &Service{
		repo:          repo,
		genresService: genresService,
		starService:   starService,
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
	return s.assemble(ctx, movie)
}

func (s *Service) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	m, err := s.repo.GetMovieByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = s.assemble(ctx, m)
	return m, err
}

func (s *Service) GetAllMoviesPaginated(ctx context.Context, searchTerm *string, sortByRating *string, starID *int, offset int, limit int) ([]*MovieDetails, int, error) {
	return s.repo.GetAllMoviesPaginated(ctx, searchTerm, sortByRating, starID, offset, limit)
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

func (s *Service) assemble(ctx context.Context, movie *MovieDetails) error {
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		movie.Genres, err = s.genresService.GetGenreByMovieID(groupCtx, movie.ID)
		return err
	})
	group.Go(func() error {
		var err error
		movie.Cast, err = s.starService.GetCastByMovieID(groupCtx, movie.ID)
		return err
	})
	return group.Wait()
}
