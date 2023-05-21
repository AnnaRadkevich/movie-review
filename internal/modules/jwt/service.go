package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Service struct {
	secret           string
	accessExpiration time.Duration
}

func NewService(secret string, accessExpiration time.Duration) *Service {
	return &Service{
		secret:           secret,
		accessExpiration: accessExpiration,
	}
}
func (s *Service) GenerateToken(id int, role string) (string, error) {
	now := time.Now()
	claims := &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   strconv.Itoa(id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiration)),
		},
		UserID: id,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}
