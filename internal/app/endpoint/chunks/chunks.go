package chunks

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"github.com/lapkomo2018/DiskordServer/pkg/hash"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
)

type Chunks struct {
	db      *gorm.DB
	discord model.DiscordService
}

func New(db *gorm.DB, discord model.DiscordService) *Chunks {
	return &Chunks{
		db:      db,
		discord: discord,
	}
}

func (c *Chunks) Upload(ctx *fiber.Ctx) error {
	var err error
	// get file from local context
	file, ok := ctx.Locals("file").(model.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	// parse body
	var body struct {
		Hash  string
		Size  uint64
		Index uint64
	}
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	// check is file need new chunks
	if err := c.db.Preload("Chunks").Find(&file).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to find file or preload chunks")
	}

	if !file.IsNeedChunk(body.Size) {
		return fiber.NewError(fiber.StatusBadRequest, "File doesnt need this chunk")
	}

	// get chunkFile
	var chunkFile *multipart.FileHeader
	if chunkFile, err = ctx.FormFile("file"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse file from body")
	}

	// checking chunkFileSize
	if uint64(chunkFile.Size) > body.Size {
		return fiber.NewError(fiber.StatusBadRequest, "ChunkFile size is too big with this chunk")
	}

	chunkFileReader, err := chunkFile.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open chunk file")
	}
	defer chunkFileReader.Close()

	// get ChunkFile bytes
	chunkBytes, err := io.ReadAll(chunkFileReader)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read chunk file")
	}

	// check Chunks hash
	var hashString string
	if hashString, err = hash.CalculateFromFile(bytes.NewReader(chunkBytes)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate hash")
	}
	if hashString != body.Hash {
		return fiber.NewError(fiber.StatusBadRequest, "Chunk hash does not match")
	}

	// upload chunkFile into discord
	var messageId string
	messageId, err = c.discord.UploadFile(chunkFile.Filename, bytes.NewReader(chunkBytes))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload chunk file to Discord")
	}

	// create chunk
	chunk := model.Chunk{
		FileId:    file.ID,
		Hash:      body.Hash,
		Size:      body.Size,
		Index:     uint(body.Index),
		MessageId: messageId,
	}

	// upload to bd
	if err := c.db.Create(&chunk).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create chunk in BD")
	}

	// response
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
