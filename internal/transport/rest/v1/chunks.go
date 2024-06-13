package v1

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"mime/multipart"
)

func (h *Handler) initChunksRouters(api fiber.Router) {
	chunks := api.Group("/chunks")
	{
		chunks.Post("/upload", skip.New(h.userIdentify, isInLocals(UserLocals)), fileOwnerCheck, h.chunkUpload)

		chunk := chunks.Group(fmt.Sprintf("/:%s<min(0)>", ChunkIndexParams), h.setChunkFromRequest)
		{
			chunk.Get("/", chunkInfo)
			chunk.Get("/download", h.chunkDownload)
		}
	}
}

func (h *Handler) chunkDownload(ctx *fiber.Ctx) error {
	chunk, ok := ctx.Locals(ChunkLocals).(core.Chunk)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse chunk")
	}

	// open reader
	fileReader, err := h.chunkService.DownloadChunk(&chunk)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download chunk file")
	}
	// preload
	return ctx.Status(fiber.StatusOK).SendStream(fileReader)
}

func chunkInfo(ctx *fiber.Ctx) error {
	chunk, ok := ctx.Locals(ChunkLocals).(core.Chunk)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse chunk")
	}

	responseChunk := struct {
		ID   uint   `json:"id"`
		Size uint64 `json:"size"`
		Hash string `json:"hash"`
	}{
		ID:   chunk.ID,
		Size: chunk.Size,
		Hash: chunk.Hash,
	}
	return ctx.Status(fiber.StatusOK).JSON(responseChunk)
}

func (h *Handler) chunkUpload(ctx *fiber.Ctx) error {
	file, ok := ctx.Locals(FileLocals).(core.File)
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

	isNeed, err := h.fileService.IsNeedChunk(&file, body.Size)
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
	if err := h.chunkService.UploadChunk(&chunk, chunkFile); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
