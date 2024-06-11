package chunk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
	"strconv"
)

type Chunk struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Chunk {
	return &Chunk{
		db: db,
	}
}

func (c *Chunk) Require(ctx *fiber.Ctx) error {
	var err error
	// get file from local storage
	file, ok := ctx.Locals("file").(model.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse BD file")
	}

	// get chunkIndex
	var chunkIndex int
	chunkIndex, err = strconv.Atoi(ctx.Params("chunkIndex"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid chunk index")
	}

	// get chunk from bd
	var chunk model.Chunk
	if err := c.db.First(&chunk, &model.Chunk{FileId: file.ID, Index: uint(chunkIndex)}).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chunk not found")
	}

	ctx.Locals("chunk", chunk)
	return ctx.Next()
}
