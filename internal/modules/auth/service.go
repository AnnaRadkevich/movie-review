package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userService *users.Service
	jwtService  *jwt.Service
}

func NewService(userService *users.Service, jwtService *jwt.Service) *Service {
	return &Service{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *Service) Register(ctx context.Context, user *users.User, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userWithPassword := &users.UserWithPassword{
		User:         user,
		PasswordHash: string(passHash),
	}

	return s.userService.Create(ctx, userWithPassword)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userService.GetExistingUserWithPassword(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("wrong email or password")
	}

	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("Generate token: %w", err)
	}
	return accessToken, nil
}
