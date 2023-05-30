package echox

import (
	"errors"
	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type HttpError struct {
	Message    string `json:"message"`
	IncidentId string `json:"incidentId,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var appError *apperrors.Error
	if !errors.As(err, &appError) {
		appError = apperrors.InternalWithoutStackTrace(err)
	}
	httpError := HttpError{
		Message:    appError.SafeError(),
		IncidentId: appError.IncidentId,
	}
	if appError.Code == apperrors.InternalCode {
		log.Printf("[ERROR] %s %s : %s\nincidentId: %s\nstack trace: %s",
			c.Request().Method, c.Request().RequestURI, err.Error(), appError.IncidentId, appError.StackTrace)
	} else {
		log.Printf("[WARN] %s %s : %s", c.Request().Method, c.Request().RequestURI, err.Error())
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
