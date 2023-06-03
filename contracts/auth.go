package contracts

type RegisterUserRequest struct {
	Username string `json:"username" validate:"min=5,max=16" `
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}
type AuthenticadedRequest[T any] struct {
	AccessToken string
	Request     T
}

func NewAuthenticated[T any](req T, accessToken string) *AuthenticadedRequest[T] {
	return &AuthenticadedRequest[T]{Request: req, AccessToken: accessToken}
}
