package v1

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"mime/multipart"
	"net/http"
)

func (h *Handler) initChunksRouters(filePublic *echo.Group, filePrivate *echo.Group) {
	chunksPublic := filePublic.Group("/chunks")
	{
		chunk := chunksPublic.Group(fmt.Sprintf("/:%s", ChunkIndexParams), h.chunkIdentify)
		{
			chunk.GET("", chunkInfo)
			chunk.GET("/download", h.chunkDownload)
		}
	}
	chunksPrivate := filePrivate.Group("/chunks")
	{
		chunksPrivate.POST("/upload", h.chunkUpload)
	}
}

func (h *Handler) chunkDownload(c echo.Context) error {
	chunk, ok := c.Get(ChunkLocals).(core.Chunk)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse chunk")
	}

	fileReader, err := h.chunkService.DownloadChunk(&chunk)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to download chunk file")
	}

	return c.Stream(http.StatusOK, "application/octet-stream", fileReader)
}

type (
	chunkInfoOutput struct {
		ID   uint   `json:"id" form:"id" query:"id"`
		Size uint64 `json:"size" form:"size" query:"size"`
		Hash string `json:"hash" form:"hash" query:"hash"`
	}
)

func chunkInfo(c echo.Context) error {
	chunk, ok := c.Get(ChunkLocals).(core.Chunk)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse chunk")
	}

	responseChunk := chunkInfoOutput{
		ID:   chunk.ID,
		Size: chunk.Size,
		Hash: chunk.Hash,
	}
	return c.JSON(http.StatusOK, responseChunk)
}

type (
	chunkUploadInput struct {
		Hash  string `form:"hash" json:"hash" binding:"required"`
		Size  uint64 `form:"size" json:"size" binding:"required"`
		Index uint   `form:"index" json:"index" binding:"required"`
	}
)

func (h *Handler) chunkUpload(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
	}

	var input chunkUploadInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse body")
	}

	isNeed, err := h.fileService.IsNeedChunk(&file, input.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !isNeed {
		return echo.NewHTTPError(http.StatusBadRequest, "File doesnt need this chunk")
	}

	var chunkFile *multipart.FileHeader
	if chunkFile, err = c.FormFile("file"); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse file from body")
	}

	if uint64(chunkFile.Size) != input.Size {
		return errors.New("file size is invalid")
	}

	chunk := core.Chunk{
		FileId: file.ID,
		Hash:   input.Hash,
		Size:   input.Size,
		Index:  input.Index,
	}
	if err := h.chunkService.UploadChunk(&chunk, chunkFile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{})
}
