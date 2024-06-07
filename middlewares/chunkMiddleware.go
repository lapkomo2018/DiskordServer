package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
	"strconv"
)

func RequireChunk(c *fiber.Ctx) error {
	var err error
	// get file from local storage
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse BD file")
	}

	// get chunkIndex
	var chunkIndex int
	chunkIndex, err = strconv.Atoi(c.Params("chunkIndex"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid chunk index")
	}

	// get chunk from bd
	var chunk models.Chunk
	if err := initializers.DB.First(&chunk, &models.Chunk{FileId: file.ID, Index: uint(chunkIndex)}).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chunk not found")
	}

	c.Locals("chunk", chunk)
	return c.Next()
}
