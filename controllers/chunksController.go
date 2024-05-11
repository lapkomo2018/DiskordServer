package controllers

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/functions"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
	"io"
	"mime/multipart"
)

func UploadChunk(c *fiber.Ctx) error {
	var err error
	// get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse file",
		})
	}

	// parse body
	var body struct {
		Hash  string
		Size  uint64
		Index uint64
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse body",
		})
	}

	// check is file need new chunks
	if err := initializers.DB.Preload("Chunks").Find(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find file or preload chunks",
		})
	}
	var totalSize uint64
	for _, chunk := range file.Chunks {
		totalSize += chunk.Size
	}
	totalSize += body.Size
	if totalSize > file.Size {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size is too big with this chunk",
		})
	}

	// get chunkFile
	var chunkFile *multipart.FileHeader
	if chunkFile, err = c.FormFile("file"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse file from body",
		})
	}

	// checking chunkFileSize
	if uint64(chunkFile.Size) > body.Size {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ChunkFile size is too big with this chunk",
		})
	}

	chunkFileReader, err := chunkFile.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to open chunk file",
		})
	}
	defer chunkFileReader.Close()

	// get ChunkFile bytes
	chunkBytes, err := io.ReadAll(chunkFileReader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read chunk file",
		})
	}

	// check Chunks hash
	var hash string
	if hash, err = functions.CalculateHashFormFile(bytes.NewReader(chunkBytes)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to calculate hash",
		})
	}
	if hash != body.Hash {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chunk hash does not match",
		})
	}

	// upload chunkFile into discord
	var message *discordgo.Message
	message, err = initializers.DiscordBot.ChannelFileSend(initializers.ChannelID, chunkFile.Filename, bytes.NewReader(chunkBytes))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to upload chunk file to Discord",
		})
	}

	// create chunk
	chunk := models.Chunk{
		FileId:    file.ID,
		Hash:      body.Hash,
		Size:      body.Size,
		Index:     uint(body.Index),
		MessageId: message.ID,
	}

	// upload to bd
	if err := initializers.DB.Create(&chunk).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload chunk in bd",
		})
	}

	// response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func DownloadChunk(c *fiber.Ctx) error {
	var err error
	// get chunk from local
	chunk, ok := c.Locals("chunk").(models.Chunk)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse chunk",
		})
	}

	// open reader
	var fileReader io.Reader
	fileReader, err = initializers.DownloadFilesFromMessage(initializers.ChannelID, chunk.MessageId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to download chunk file",
		})
	}

	// preload
	return c.Status(fiber.StatusOK).SendStream(fileReader)
}

func GetChunkInfo(c *fiber.Ctx) error {
	// get chunk from local
	chunk, ok := c.Locals("chunk").(models.Chunk)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse chunk",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"chunk": chunk,
	})
}
