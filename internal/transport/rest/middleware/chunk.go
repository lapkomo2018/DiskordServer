package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"strconv"
)

type ChunkService interface {
	First(chunk *core.Chunk, cond ...interface{}) error
}

type ChunkMiddleware struct {
	service ChunkService
}

func NewChunkMiddleware(s ChunkService) *ChunkMiddleware {
	return &ChunkMiddleware{
		service: s,
	}
}

func (c *ChunkMiddleware) Require(ctx *fiber.Ctx) error {
	var err error
	// get file from local storage
	file, ok := ctx.Locals("file").(core.File)
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
	var chunk core.Chunk
	if err := c.service.First(&chunk, &core.Chunk{FileId: file.ID, Index: uint(chunkIndex)}); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chunk not found")
	}

	ctx.Locals("chunk", chunk)
	return ctx.Next()
}
