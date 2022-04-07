package error

import (
	"github.com/labstack/echo/v4"
)

// Error struct
type Error struct {
	Code    int         `json:"code" example:"409"`
	Message string      `json:"message" example:"Some Error"`
	Details interface{} `json:"details,omitempty" example:"[]"`
}

// HTTPError Construct an HTTP error
func HttpError(status int, msg string) *echo.HTTPError {
	e := Error{
		Code:    status,
		Message: msg,
	}
	return echo.NewHTTPError(status, e)
}

// DetailedHTTPError Construct an HTTP error
func DetailedHttpError(status int, msg string, details interface{}) *echo.HTTPError {
	e := Error{
		Code:    status,
		Message: msg,
		Details: details,
	}
	return echo.NewHTTPError(status, e)
}
