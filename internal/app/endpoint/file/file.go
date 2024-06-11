package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
)

type File struct {
	db *gorm.DB
}

func New(db *gorm.DB) *File {
	return &File{
		db: db,
	}
}

func (f *File) Download(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(model.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	//preload pieces
	if err := f.db.Preload("Chunks").Find(&file).Error; err != nil {
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

func (f *File) Info(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(model.File)
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

func (f *File) ChangePrivacy(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(model.File)
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
	if err := f.db.Save(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to patch file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (f *File) Delete(c *fiber.Ctx) error {
	//get file from local context
	file, ok := c.Locals("file").(model.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	if err := f.db.Delete(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
