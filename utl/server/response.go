package server

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

var (
	// Error message
	ErrSizeNotFound        = errors.New("request size not found")
	ErrSizeMenuNotFound    = errors.New("request size menu not found")
	ErrSizeProfileNotFound = errors.New("request size profile not found")
	ErrSizeProductNotFound = errors.New("request size product not found")

	ErrInvalidUsernamePassword = errors.New("Invalid username/password")
)

type ResponseMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NewResponse struct {
	Meta ResponseMeta `json:"meta"`
	Data interface{}  `json:"data,omitempty"`
}

func ResponseOK(c echo.Context, msg string, r interface{}) error {
	return c.JSON(http.StatusOK, NewResponse{
		Meta: ResponseMeta{
			Code:    http.StatusOK,
			Message: msg,
		},
		Data: r,
	})
}

func ResponseFail(c echo.Context, err error) error {
	return c.JSON(http.StatusOK, NewResponse{
		Meta: ResponseMeta{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		},
	})
}

func ResponseUnauthorized(c echo.Context, err error) error {
	return c.JSON(http.StatusOK, NewResponse{
		Meta: ResponseMeta{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		},
	})
}

func NotFound(c echo.Context, err error) error {
	return c.JSON(http.StatusNotFound, NewResponse{
		Meta: ResponseMeta{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		},
	})
}
