package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
)

type FileService interface {
	Create(file *core.File) error
	Save(file *core.File) error
	Delete(file *core.File) error
	IsNeedChunk(file *core.File, chunkSize uint64) (bool, error)
	LoadChunks(file *core.File) error
}

type FileHandler struct {
	service FileService
}

func NewFileHandler(s FileService) *FileHandler {
	return &FileHandler{
		service: s,
	}
}

func (f *FileHandler) Download(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	//preload pieces
	if err := f.service.LoadChunks(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load user files")
	}

	// validate file
	if !file.Validate() {
		return fiber.NewError(fiber.StatusInternalServerError, "File is corrupted")
	}

	type resChunk struct {
		Index uint `json:"index"`
	}
	var resChunks []resChunk
	for _, chunk := range file.Chunks {
		resChunks = append(resChunks, resChunk{
			Index: chunk.Index,
		})
	}

	return c.Status(fiber.StatusOK).JSON(resChunks)
}

func (f *FileHandler) Info(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(core.File)
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

func (f *FileHandler) ChangePrivacy(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(core.File)
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
	if err := f.service.Save(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to patch file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (f *FileHandler) Delete(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	if err := f.service.Delete(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (f *FileHandler) Upload(c *fiber.Ctx) error {
	//get user from local context
	user, ok := c.Locals("user").(core.User)
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
	file := core.File{
		UserId:    user.ID,
		Name:      body.Name,
		Hash:      body.Hash,
		Size:      body.Size,
		IsPublic:  body.IsPublic,
		NumChunks: body.NumChunks,
		ChunkSize: body.ChunkSize,
	}

	if err := f.service.Create(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create file")
	}

	//response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileId": file.ID,
	})

}
