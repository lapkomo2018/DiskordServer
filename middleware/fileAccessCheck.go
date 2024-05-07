package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
)

func FileAccessCheck(c *fiber.Ctx) error {
	// get user
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	// parse file id from body
	var body struct {
		FileID string
	}
	if c.BodyParser(&body) != nil {
		body.FileID = c.FormValue("fileID")
		if body.FileID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to read body",
			})
		}
	}

	// get file from bd
	var file models.File
	if result := initializers.DB.First(&file, body.FileID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// check file is public
	if !file.IsPublic {
		// check file owner
		if file.UserID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not allowed to access this file",
			})
		}
	}

	// putting file into local
	c.Locals("file", file)
	return c.Next()
}
