package rest

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"net/http"
)

func ErrorHandler(err error, c echo.Context) {
	// Status code defaults to 500
	code := http.StatusInternalServerError
	message := http.StatusText(code)

	// Check if it's an HTTP error
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	// Return JSON response with error message
	_ = c.JSON(code, core.ErrorResponse{
		Error: message,
	})
}
