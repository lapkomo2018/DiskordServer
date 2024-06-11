package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
	"strconv"
)

type File struct {
	db *gorm.DB
}

func New(db *gorm.DB) *File {
	return &File{
		db: db,
	}
}

func (f *File) Require(c *fiber.Ctx) error {
	var err error
	// get fileId
	var fileId int
	fileId, err = strconv.Atoi(c.Params("fileId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file ID")
	}

	// get file from bd
	var file model.File
	if err := f.db.First(&file, fileId).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	// putting file into local
	c.Locals("file", file)
	return c.Next()
}

func (f *File) OwnerCheck(c *fiber.Ctx) error {
	// get user
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	// get file
	file, ok := c.Locals("file").(model.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	// check file owner
	if file.UserId != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "You are not the owner of the file")
	}
	return c.Next()
}

func (f *File) IsPublic(c *fiber.Ctx) bool {
	// get file
	file, ok := c.Locals("file").(model.File)
	if !ok {
		return false
	}
	return file.IsPublic
}
