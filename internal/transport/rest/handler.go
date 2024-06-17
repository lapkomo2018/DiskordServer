package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"net/http"
)

func ErrorHandler(err error, c echo.Context) {
	// Status code defaults to 500
	code := http.StatusInternalServerError
	message := err.Error()

	// Check if it's an HTTP error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Return JSON response with error message
	_ = c.JSON(code, core.ErrorResponse{
		Error: message,
	})
}
