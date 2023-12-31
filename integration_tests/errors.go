package tests

import (
	"net/http"
	"testing"

	"github.com/RadkevichAnn/movie-reviews/client"
	"github.com/RadkevichAnn/movie-reviews/internal/apperrors"
	"github.com/stretchr/testify/require"
)

func requireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireApiError(t, err, http.StatusNotFound, msg)
}

func requireUnauthorizedError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusUnauthorized, msg)
}

func requireForbiddenError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusForbidden, msg)
}

func requireBadRequestError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusBadRequest, msg)
}

func requireVersionMismatchError(t *testing.T, err error, subject, key string, value any, version int) {
	msg := apperrors.VersionMismatch(subject, key, value, version).Error()
	requireApiError(t, err, http.StatusConflict, msg)
}

func requireAlreadyExistError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.AlreadyExists(subject, key, value).Error()
	requireApiError(t, err, http.StatusConflict, msg)
}

func requireApiError(t *testing.T, err error, statusCode int, msg string) {
	cerr, ok := err.(*client.Error)
	require.True(t, ok, "expected client.Error")
	require.Equal(t, statusCode, cerr.Code)
	require.Contains(t, cerr.Message, msg)
}
