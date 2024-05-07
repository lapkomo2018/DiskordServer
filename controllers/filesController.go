package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
)

func GetUserFilesList(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	if err := initializers.DB.Preload("Files").Preload("Files.Chunks").Find(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load user files",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files": user.Files,
	})
}

func UploadFile(c *fiber.Ctx) error {
	//get user from local context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	//get file info from body
	var body struct {
		Name      string
		Size      uint64
		Hash      string
		IsPublic  bool
		NumChunks uint
		ChunkSize uint64
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse body",
		})
	}

	//add file to bd
	file := models.File{
		UserID:    user.ID,
		Name:      body.Name,
		Hash:      body.Hash,
		Size:      body.Size,
		IsPublic:  body.IsPublic,
		NumChunks: body.NumChunks,
		ChunkSize: body.ChunkSize,
	}
	result := initializers.DB.Create(&file)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	//response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileID": file.ID,
	})

}

func DownloadFile(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse file",
		})
	}

	//preload pieces
	if err := initializers.DB.Preload("Chunks").Find(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load user files",
		})
	}

	// validate file
	if file.NumChunks != uint(len(file.Chunks)) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "File is corrupted",
		})
	}
	var totalSize uint64
	for _, chunk := range file.Chunks {
		totalSize += chunk.Size
	}
	if totalSize != file.Size {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "File is corrupted",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name":   file.Name,
		"size":   file.Size,
		"hash":   file.Hash,
		"chunks": file.Chunks,
	})
}
