package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
	"strconv"
)

func RequireFile(c *fiber.Ctx) error {
	var err error
	// get fileId
	var fileId int
	fileId, err = strconv.Atoi(c.Params("fileId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	// get file from bd
	var file models.File
	if err := initializers.DB.First(&file, fileId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// putting file into local
	c.Locals("file", file)
	return c.Next()
}

func FileOwnerCheck(c *fiber.Ctx) error {
	// get user
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	// get file
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse file",
		})
	}

	// check file owner
	if file.UserId != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not the owner of the file",
		})
	}
	return c.Next()
}

func FileIsPublic(c *fiber.Ctx) bool {
	// get file
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return false
	}
	return file.IsPublic
}
