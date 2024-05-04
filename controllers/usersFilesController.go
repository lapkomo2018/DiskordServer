package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
)

func loadUserFilesAndPieces(user *models.User) error {
	return initializers.DB.Preload("Files").Preload("Files.Pieces").Find(user).Error
}

func GetUserFilesList(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	if err := loadUserFilesAndPieces(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load user files",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files": user.Files,
	})
}

func DownloadFile(c *fiber.Ctx) error {
	//get user from local context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	//preload files and pieces
	if err := loadUserFilesAndPieces(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load user files",
		})
	}

	//

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files": user.Files,
	})
}
