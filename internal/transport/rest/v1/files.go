package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/lapkomo2018/DiskordServer/internal/core"
)

func (h *Handler) initFilesRouters(api fiber.Router) {
	files := api.Group("/files")
	{
		files.Post("/upload", h.userIdentify, h.fileUpload)

		file := files.Group(fmt.Sprintf("/:%s<min(1)>", FileIdParams), h.setFileFromRequest, skip.New(h.userIdentify, fileIsPublic), skip.New(fileOwnerCheck, fileIsPublic))
		{
			file.Get("/", fileInfo)
			file.Get("/download", h.fileChunks)
			file.Patch("/privacy", skip.New(h.userIdentify, isInLocals(UserLocals)), fileOwnerCheck, h.fileChangePrivacy)
			file.Delete("/", skip.New(h.userIdentify, isInLocals(UserLocals)), fileOwnerCheck, h.fileDelete)

			h.initChunksRouters(file)
		}
	}
}

func (h *Handler) fileChunks(c *fiber.Ctx) error {
	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	//preload pieces
	if err := h.fileService.LoadChunks(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load user files")
	}

	if !file.Validate() {
		return fiber.NewError(fiber.StatusInternalServerError, "File is corrupted")
	}

	type resChunk struct {
		Index uint `json:"index"`
	}
	var resChunks []resChunk
	for _, chunk := range file.Chunks {
		resChunks = append(resChunks, resChunk{
			Index: chunk.Index,
		})
	}

	return c.Status(fiber.StatusOK).JSON(resChunks)
}

func fileInfo(c *fiber.Ctx) error {
	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	responseFile := struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
		Size uint64 `json:"size"`
		Hash string `json:"hash"`
	}{
		Id:   file.ID,
		Name: file.Name,
		Size: file.Size,
		Hash: file.Hash,
	}
	return c.Status(fiber.StatusOK).JSON(responseFile)
}

func (h *Handler) fileChangePrivacy(c *fiber.Ctx) error {
	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	var body struct {
		IsPublic bool
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	file.IsPublic = body.IsPublic
	if err := h.fileService.Save(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to patch file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *Handler) fileDelete(c *fiber.Ctx) error {
	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	if err := h.fileService.Delete(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *Handler) fileUpload(c *fiber.Ctx) error {
	user, ok := c.Locals(UserLocals).(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	//get file info from body
	var body struct {
		Name      string
		Size      uint64
		Hash      string
		IsPublic  bool
		NumChunks uint
		ChunkSize uint64
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	//add file to bd
	file := core.File{
		UserId:    user.ID,
		Name:      body.Name,
		Hash:      body.Hash,
		Size:      body.Size,
		IsPublic:  body.IsPublic,
		NumChunks: body.NumChunks,
		ChunkSize: body.ChunkSize,
	}

	if err := h.fileService.Create(&file); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create file")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileId": file.ID,
	})
}
