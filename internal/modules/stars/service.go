package stars

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateStar(ctx context.Context, star *Star) error {
	return s.repo.CreateStar(ctx, star)
}

func (s *Service) GetStarByID(ctx context.Context, id int) (*Star, error) {
	return s.repo.GetStarByID(ctx, id)
}
