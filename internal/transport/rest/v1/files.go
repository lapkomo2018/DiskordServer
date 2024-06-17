package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"net/http"
)

func (h *Handler) initFilesRouters(api *echo.Group) {
	files := api.Group("/files")
	{
		files.POST("/upload", h.fileUpload, h.userIdentify)

		file := files.Group(fmt.Sprintf("/:%s", FileIdParams), h.fileIdentify)
		{
			filePublic := file.Group("", h.fileAccessCheck(true))
			{
				filePublic.GET("", fileInfo)
				filePublic.GET("/download", h.fileChunks)
			}

			filePrivate := file.Group("", h.fileAccessCheck(false))
			{
				filePrivate.DELETE("", h.fileDelete)
				filePrivate.PATCH("/privacy", h.fileChangePrivacy)
			}

			h.initChunksRouters(filePublic, filePrivate)
		}
	}
}

type (
	fileChunksOutputChunk struct {
		Index uint `json:"index" form:"index" query:"index"`
	}
)

func (h *Handler) fileChunks(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
	}

	if err := h.fileService.LoadChunks(&file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load user files")
	}

	if !file.Validate() {
		return echo.NewHTTPError(http.StatusInternalServerError, "File is corrupted")
	}

	var resChunks []fileChunksOutputChunk
	for _, chunk := range file.Chunks {
		resChunks = append(resChunks, fileChunksOutputChunk{
			Index: chunk.Index,
		})
	}

	return c.JSON(http.StatusOK, resChunks)
}

type fileInfoOutput struct {
	Id       uint   `json:"id" form:"id" query:"id"`
	Name     string `json:"name" form:"name" query:"name"`
	Size     uint64 `json:"size" form:"size" query:"size"`
	Hash     string `json:"hash" form:"hash" query:"hash"`
	IsPublic bool   `json:"isPublic" form:"isPublic" query:"isPublic"`
}

func fileInfo(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
	}

	responseFile := fileInfoOutput{
		Id:       file.ID,
		Name:     file.Name,
		Size:     file.Size,
		Hash:     file.Hash,
		IsPublic: file.IsPublic,
	}
	return c.JSON(http.StatusOK, responseFile)
}

type (
	fileChangePrivacyInput struct {
		IsPublic bool `json:"isPublic" form:"isPublic" query:"isPublic" binding:"required"`
	}
)

func (h *Handler) fileChangePrivacy(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
	}

	var input fileChangePrivacyInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse body")
	}

	file.IsPublic = input.IsPublic
	if err := h.fileService.Save(&file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to patch file")
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func (h *Handler) fileDelete(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
	}

	if err := h.fileService.Delete(&file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete file")
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

type (
	fileUploadInput struct {
		Name      string `json:"name" form:"name" query:"name" binding:"required"`
		Size      uint64 `json:"size" form:"size" query:"size" binding:"required"`
		Hash      string `json:"hash" form:"hash" query:"hash" binding:"required"`
		IsPublic  bool   `json:"isPublic" form:"isPublic" query:"isPublic" binding:"required"`
		NumChunks uint   `json:"numChunks" form:"numChunks" query:"numChunks" binding:"required"`
		ChunkSize uint64 `json:"chunkSize" form:"chunkSize" query:"chunkSize" binding:"required"`
	}
	fileUploadOutput struct {
		FileId uint `json:"fileId" form:"fileId" query:"fileId"`
	}
)

func (h *Handler) fileUpload(c echo.Context) error {
	user, ok := c.Get(UserLocals).(core.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse user")
	}

	//get file info from body
	var input fileUploadInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse body")
	}

	//add file to bd
	file := core.File{
		UserId:    user.ID,
		Name:      input.Name,
		Hash:      input.Hash,
		Size:      input.Size,
		IsPublic:  input.IsPublic,
		NumChunks: input.NumChunks,
		ChunkSize: input.ChunkSize,
	}
	if err := h.fileService.Create(&file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file")
	}

	output := fileUploadOutput{
		FileId: file.ID,
	}

	return c.JSON(http.StatusOK, output)
}
