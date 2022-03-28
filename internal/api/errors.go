package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/romarq/visualtez-testing/internal/logger"
)

// Error struct
type Error struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"Some Error"`
}

// HTTPError Construct an HTTP error
func HTTPError(ctx echo.Context, status int, msg string) error {
	e := Error{
		Code:    status,
		Message: msg,
	}
	return echo.NewHTTPError(status, e)
}

// HandleError Construct an HTTP error
func HandleError(ctx echo.Context, err error) error {
	status := http.StatusInternalServerError
	msg := "Unknown Error"

	if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
		msg = "Not Found."
	}

	logger.Debug("HTTP error: %v", err)
	return echo.NewHTTPError(status, Error{
		Code:    status,
		Message: msg,
	})
}
