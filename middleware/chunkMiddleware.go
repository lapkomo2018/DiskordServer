package middleware

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse bd file",
		})
	}

	// get chunkIndex
	var chunkIndex int
	chunkIndex, err = strconv.Atoi(c.Params("chunkIndex"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chunk index",
		})
	}

	// get chunk from bd
	var chunk models.Chunk
	if err := initializers.DB.Where(&models.Chunk{Index: uint(chunkIndex), File: file}).First(&chunk).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chunk not found",
		})
	}

	c.Locals("chunk", chunk)
	return c.Next()
}
