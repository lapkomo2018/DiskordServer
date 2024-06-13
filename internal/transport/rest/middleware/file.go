package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"strconv"
)

type FileService interface {
	First(file *core.File, cond ...interface{}) error
}

type FileMiddleware struct {
	service FileService
}

func NewFileMiddleware(s FileService) *FileMiddleware {
	return &FileMiddleware{
		service: s,
	}
}

func (f *FileMiddleware) Require(c *fiber.Ctx) error {
	var err error
	// get fileId
	var fileId int
	fileId, err = strconv.Atoi(c.Params("fileId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file ID")
	}

	// get file from bd
	var file core.File
	if err := f.service.First(&file, fileId); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	// putting file into local
	c.Locals("file", file)
	return c.Next()
}

func (f *FileMiddleware) OwnerCheck(c *fiber.Ctx) error {
	// get user
	user, ok := c.Locals("user").(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	// get file
	file, ok := c.Locals("file").(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	// check file owner
	if file.UserId != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "You are not the owner of the file")
	}
	return c.Next()
}

func (f *FileMiddleware) IsPublic(c *fiber.Ctx) bool {
	// get file
	file, ok := c.Locals("file").(core.File)
	if !ok {
		return false
	}
	return file.IsPublic
}
