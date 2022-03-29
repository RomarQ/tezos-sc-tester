package api

import (
	"github.com/labstack/echo/v4"
)

// Error struct
type Error struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"Some Error"`
}

// HTTPError Construct an HTTP error
func HTTPError(status int, msg string) *echo.HTTPError {
	e := Error{
		Code:    status,
		Message: msg,
	}
	return echo.NewHTTPError(status, e)
}
