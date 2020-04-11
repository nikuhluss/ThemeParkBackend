package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// ReviewHandler handles HTTP requests for review jobs.
type ReviewHandler struct {
	reviewUsecase usecases.ReviewUsecase
}

// NewReviewHandler returns a new ReviewHanler instance.
func NewReviewHandler(reviewUsecase usecases.ReviewUsecase) *ReviewHandler {
	return &ReviewHandler{
		reviewUsecase,
	}
}

// Bind sets up the routes for the handler.
func (rh *ReviewHandler) Bind(e *echo.Echo) error {
	e.GET("/reviews", rh.Fetch)
	e.POST("/reviews", rh.Store)
	e.GET("/reviews/:reviewID", rh.GetByID)
	e.PUT("/reviews/:reviewID", rh.Update)
	e.DELETE("/reviews/:reviewID", rh.Delete)
	e.GET("/rides/:rideID/reviews", rh.FetchForRide)
	return nil
}

// GetByID gets a specific review.
func (rh *ReviewHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	reviewID := c.Param("reviewID")

	review, err := rh.reviewUsecase.GetByID(ctx, reviewID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, review, Indent)
}

// Fetch fetches all reviews.
func (rh *ReviewHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	reviews, err := rh.reviewUsecase.Fetch(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, reviews, Indent)
}

// FetchForRide fetches all reviews for the given ride.
func (rh *ReviewHandler) FetchForRide(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("rideID")

	reviews, err := rh.reviewUsecase.FetchForRide(ctx, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, reviews, Indent)
}

// Store creates a new review.
func (rh *ReviewHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	review := &models.Review{}

	err := c.Bind(review)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = rh.reviewUsecase.Store(ctx, review)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, review, Indent)
}

// Update updates a specific review.
func (rh *ReviewHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	reviewID := c.Param("reviewID")

	review := &models.Review{}
	review.ID = reviewID

	err := c.Bind(review)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = rh.reviewUsecase.Update(ctx, review)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, review, Indent)
}

// Delete deletes a specific review.
func (rh *ReviewHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	reviewID := c.Param("reviewID")

	err := rh.reviewUsecase.Delete(ctx, reviewID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, "", Indent)
}
