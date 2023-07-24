package reviews

import (
	"errors"
	"net/http"

	"golang.org/x/sync/singleflight"

	"github.com/RadkevichAnn/movie-reviews/contracts"
	"github.com/RadkevichAnn/movie-reviews/internal/apperrors"
	"github.com/RadkevichAnn/movie-reviews/internal/config"
	"github.com/RadkevichAnn/movie-reviews/internal/echox"
	"github.com/RadkevichAnn/movie-reviews/internal/pagination"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig config.PaginationConfig
	reqGroup         singleflight.Group
}

func NewHandler(service *Service, paginationConfig config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) CreateReview(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.CreateReviewRequest](c)
	if err != nil {
		return err
	}
	review := &Review{
		MovieID: req.MovieID,
		UserID:  req.UserID,
		Rating:  req.Rating,
		Title:   req.Title,
		Content: req.Content,
	}
	err = h.service.CreateReview(c.Request().Context(), review)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, review)
}

func (h *Handler) GetReviewByID(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetReviewRequest](c)
		if err != nil {
			return nil, err
		}
		review, err := h.service.GetReviewByID(c.Request().Context(), req.ReviewID)
		if err != nil {
			return nil, err
		}
		return review, nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllReviewsPaginated(c echo.Context) error {
	res, err, _ := h.reqGroup.Do(c.Request().RequestURI, func() (any, error) {
		req, err := echox.BindAndValidate[contracts.GetReviewsRequest](c)
		if err != nil {
			return nil, err
		}
		if req.MovieID == nil && req.UserID == nil {
			return nil, apperrors.BadRequest(errors.New("either movie_id or user_id must be provided"))
		}
		pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
		offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
		reviews, total, err := h.service.GetAllReviewsPaginated(c.Request().Context(), req.MovieID, req.UserID, offset, limit)
		if err != nil {
			return nil, err
		}
		return pagination.Response(&req.PaginatedRequest, total, reviews), nil
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateReview(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.UpdateReviewRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.UpdateReview(c.Request().Context(), req.ReviewID, req.UserID, req.Title, req.Content, req.Rating); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteReview(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.DeleteReviewRequest](c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteReview(c.Request().Context(), req.ReviewID, req.UserID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
