package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"io"
	"mime/multipart"
)

type ChunkService interface {
	DownloadChunk(chunk core.Chunk) (io.Reader, error)
	UploadChunk(chunk *core.Chunk, file *multipart.FileHeader) error
}

type ChunkHandler struct {
	chunkService ChunkService
	fileService  FileService
}

func NewChunkHandler(chunkService ChunkService, fileService FileService) *ChunkHandler {
	return &ChunkHandler{
		chunkService: chunkService,
		fileService:  fileService,
	}
}

func (c *ChunkHandler) Download(ctx *fiber.Ctx) error {
	var err error
	// get chunk from local
	chunk, ok := ctx.Locals("chunk").(core.Chunk)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse chunk")
	}

	// open reader
	var fileReader io.Reader
	fileReader, err = c.chunkService.DownloadChunk(chunk)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download chunk file")
	}
	// preload
	return ctx.Status(fiber.StatusOK).SendStream(fileReader)
}

func (c *ChunkHandler) Info(ctx *fiber.Ctx) error {
	// get chunk from local
	chunk, ok := ctx.Locals("chunk").(core.Chunk)
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

func (c *ChunkHandler) Upload(ctx *fiber.Ctx) error {
	// get file from local context
	file, ok := ctx.Locals("file").(core.File)
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

	isNeed, err := c.fileService.IsNeedChunk(&file, body.Size)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !isNeed {
		return fiber.NewError(fiber.StatusBadRequest, "File doesnt need this chunk")
	}

	// get chunkFile
	var chunkFile *multipart.FileHeader
	if chunkFile, err = ctx.FormFile("file"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse file from body")
	}

	if uint64(chunkFile.Size) != body.Size {
		return errors.New("file size is invalid")
	}

	// create chunk
	chunk := core.Chunk{
		FileId: file.ID,
		Hash:   body.Hash,
		Size:   body.Size,
		Index:  uint(body.Index),
	}

	if err := c.chunkService.UploadChunk(&chunk, chunkFile); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// response
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
