package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Handle(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		code = fiberError.Code
	}

	return c.Status(code).JSON(core.ErrorResponse{
		Error: err.Error(),
	})
}
