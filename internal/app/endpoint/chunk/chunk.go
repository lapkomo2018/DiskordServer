package chunk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
	"io"
)

type Chunk struct {
	db      *gorm.DB
	discord model.DiscordService
}

func New(db *gorm.DB, discord model.DiscordService) *Chunk {
	return &Chunk{
		db:      db,
		discord: discord,
	}
}

func (c *Chunk) Download(ctx *fiber.Ctx) error {
	var err error
	// get chunk from local
	chunk, ok := ctx.Locals("chunk").(model.Chunk)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse chunk")
	}

	// open reader
	var fileReader io.Reader
	fileReader, err = c.discord.DownloadFileFromMessage(chunk.MessageId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download chunk file")
	}
	// preload
	return ctx.Status(fiber.StatusOK).SendStream(fileReader)
}

func (c *Chunk) Info(ctx *fiber.Ctx) error {
	// get chunk from local
	chunk, ok := ctx.Locals("chunk").(model.Chunk)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse chunk")
	}

	responseChunk := struct {
		ID   uint
		Size uint64
		Hash string
	}{
		ID:   chunk.ID,
		Size: chunk.Size,
		Hash: chunk.Hash,
	}
	return ctx.Status(fiber.StatusOK).JSON(responseChunk)
}
