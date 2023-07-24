package users

import (
	"context"

	"github.com/RadkevichAnn/movie-reviews/internal/log"
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
	if err := s.repo.Delete(ctx, userId); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user deleted", "userId", userId)
	return nil
}

func (s *Service) Get(ctx context.Context, userId int) (user *User, err error) {
	return s.repo.GetUserById(ctx, userId)
}

func (s *Service) UpdateBio(ctx context.Context, userId int, bio string) error {
	if err := s.repo.UpdateBio(ctx, userId, bio); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user bio updated", "userId", userId, "bio", bio)
	return nil
}

func (s *Service) UpdateRole(ctx context.Context, userId int, role string) error {
	if err := s.repo.UpdateRole(ctx, userId, role); err != nil {
		return err
	}
	log.FromContext(ctx).Info("user role updated", "userId", userId, "role", role)
	return nil
}

func (s *Service) GetExistingUserByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetExistingUserByUsername(ctx, username)
}
