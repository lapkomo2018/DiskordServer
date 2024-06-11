package error

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
)

type Error struct {
}

func New() *Error {
	return &Error{}
}

func (e *Error) Handle(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		code = fiberError.Code
	}

	return c.Status(code).JSON(model.ErrorResponse{
		Error: err.Error(),
	})
}
