package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetExistingUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetExistingUserWithPassword(ctx, email)
}

func (s *Service) Delete(ctx context.Context, userId int) error {
	return s.repo.Delete(ctx, userId)
}

func (s *Service) Get(ctx context.Context, userId int) (user *User, err error) {
	return s.repo.GetUserById(ctx, userId)
}

func (s *Service) UpdateBio(ctx context.Context, userId int, bio string) error {
	return s.repo.UpdateBio(ctx, userId, bio)
}
func (s *Service) UpdateRole(ctx context.Context, userId int, role string) error {
	return s.repo.UpdateRole(ctx, userId, role)
}
