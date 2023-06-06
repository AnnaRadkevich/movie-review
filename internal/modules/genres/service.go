package genres

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllGenres(ctx context.Context) ([]*Genre, error) {
	return s.repo.GetAllGenres(ctx)
}

func (s *Service) GetGenreById(ctx context.Context, id int) (*Genre, error) {
	return s.repo.GetGenreById(ctx, id)
}

func (s *Service) CreateGenre(ctx context.Context, name string) (*Genre, error) {
	return s.repo.CreateGenre(ctx, name)
}

func (s *Service) UpdateGenre(ctx context.Context, id int, name string) error {
	return s.repo.UpdateGenre(ctx, id, name)
}

func (s *Service) DeleteGenre(ctx context.Context, id int) error {
	return s.repo.DeleteGenre(ctx, id)
}
