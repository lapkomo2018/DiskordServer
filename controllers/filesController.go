package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
)

func GetUserFiles(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user",
		})
	}

	if err := initializers.DB.Preload("Files").Find(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load user files",
		})
	}

	var fileIds []uint
	for _, file := range user.Files {
		fileIds = append(fileIds, file.ID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileIds": fileIds,
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
		UserId:    user.ID,
		Name:      body.Name,
		Hash:      body.Hash,
		Size:      body.Size,
		IsPublic:  body.IsPublic,
		NumChunks: body.NumChunks,
		ChunkSize: body.ChunkSize,
	}

	if err := initializers.DB.Create(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileId": file.ID,
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

	var chunkIndexes []uint
	for _, chunk := range file.Chunks {
		chunkIndexes = append(chunkIndexes, chunk.Index)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"chunks": chunkIndexes,
	})
}

func GetFileInfo(c *fiber.Ctx) error {
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

	responseFile := struct {
		id   uint
		name string
		size uint64
		hash string
	}{
		id:   file.ID,
		name: file.Name,
		size: file.Size,
		hash: file.Hash,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"file": responseFile,
	})
}
