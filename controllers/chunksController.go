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
	"strconv"
)

func UploadChunk(c *fiber.Ctx) error {
	// get file from local context
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse bd file",
		})
	}

	// parse body
	var body struct {
		Hash  string
		Size  uint64
		Index uint64
	}
	var err error
	body.Hash = c.FormValue("hash")
	if body.Hash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hash is required",
		})
	}
	body.Size, err = strconv.ParseUint(c.FormValue("size"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse size",
		})
	}
	body.Index, err = strconv.ParseUint(c.FormValue("index"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse index",
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

	// create file chunk
	chunk := models.Chunk{
		FileID:    file.ID,
		Hash:      body.Hash,
		Size:      body.Size,
		Index:     uint(body.Index),
		MessageID: message.ID,
	}

	// upload to bd
	result := initializers.DB.Create(&chunk)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload chunk in bd",
		})
	}

	// response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func DownloadChunk(c *fiber.Ctx) error {
	// get file from local storage
	file, ok := c.Locals("file").(models.File)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse bd file",
		})
	}

	// preload chunks
	if result := initializers.DB.Preload("Chunks").First(&file); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find file or preload chunks",
		})
	}

	// parse body
	var body struct {
		ChunkID uint
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse body",
		})
	}

	// get chunk
	var chunk *models.Chunk
	for _, c := range file.Chunks {
		if c.ID == body.ChunkID {
			chunk = &c
			break
		}
	}
	if chunk == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chunk not found",
		})
	}
	var fileReader io.Reader
	var err error
	fileReader, err = initializers.DownloadFilesFromMessage(initializers.ChannelID, chunk.MessageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to download chunk file",
		})
	}

	// preload
	return c.Status(fiber.StatusOK).SendStream(fileReader)
}
