package echox

import (
	"errors"
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/log"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var appError *apperrors.Error
	if !errors.As(err, &appError) {
		appError = apperrors.InternalWithoutStackTrace(err)
	}
	httpError := contracts.HttpError{
		Message:    appError.SafeError(),
		IncidentId: appError.IncidentId,
	}
	logger := log.FromContext(c.Request().Context())
	if appError.Code == apperrors.InternalCode {
		logger.Error("server error", "message", err.Error(), "incidientId", appError.IncidentId, "stackTrace", appError.StackTrace)
	} else {
		logger.Warn("client error", "message", err.Error())
	}
	if err = c.JSON(toHttpStatus(appError.Code), httpError); err != nil {
		c.Logger().Error(err)
	}
}

func toHttpStatus(code apperrors.Code) int {
	switch code {
	case apperrors.InternalCode:
		return http.StatusInternalServerError
	case apperrors.BadRequestCode:
		return http.StatusBadRequest
	case apperrors.NotFoundCode:
		return http.StatusNotFound
	case apperrors.AlreadyExistsCode:
		return http.StatusConflict
	case apperrors.UnauthorizedCode:
		return http.StatusUnauthorized
	case apperrors.ForbiddenCode:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError

	}
}
