package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
)

func UploadFile(c *fiber.Ctx) error {
	//get user from local context
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
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
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create file")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	//preload pieces
	if err := initializers.DB.Preload("Chunks").Find(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load user files")
	}

	// validate file
	if file.NumChunks != uint(len(file.Chunks)) {
		return fiber.NewError(fiber.StatusInternalServerError, "File is corrupted")
	}
	var totalSize uint64
	for _, chunk := range file.Chunks {
		totalSize += chunk.Size
	}
	if totalSize != file.Size {
		return fiber.NewError(fiber.StatusInternalServerError, "File is corrupted")
	}

	type resChunk struct {
		Index uint `json:"index"`
	}
	var chunks []resChunk
	for _, chunk := range file.Chunks {
		chunks = append(chunks, resChunk{
			Index: chunk.Index,
		})
	}

	return c.Status(fiber.StatusOK).JSON(chunks)
}

func GetFileInfo(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	responseFile := struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
		Size uint64 `json:"size"`
		Hash string `json:"hash"`
	}{
		Id:   file.ID,
		Name: file.Name,
		Size: file.Size,
		Hash: file.Hash,
	}
	return c.Status(fiber.StatusOK).JSON(responseFile)
}

func ChangeFilePrivacy(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	var body struct {
		IsPublic bool
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	file.IsPublic = body.IsPublic
	if err := initializers.DB.Save(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to patch file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func DeleteFile(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	if err := initializers.DB.Delete(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
