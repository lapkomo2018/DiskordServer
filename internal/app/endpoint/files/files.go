package files

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
)

type Files struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Files {
	return &Files{
		db: db,
	}
}

func (f *Files) Upload(c *fiber.Ctx) error {
	//get user from local context
	user, ok := c.Locals("user").(model.User)
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
	file := model.File{
		UserId:    user.ID,
		Name:      body.Name,
		Hash:      body.Hash,
		Size:      body.Size,
		IsPublic:  body.IsPublic,
		NumChunks: body.NumChunks,
		ChunkSize: body.ChunkSize,
	}

	if err := f.db.Create(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create file")
	}

	//response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileId": file.ID,
	})

}
