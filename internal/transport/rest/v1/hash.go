package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/pkg/hash"
	"mime/multipart"
	"net/http"
)

func (h *Handler) initHashRouters(api *echo.Group) {
	hashRouter := api.Group("/hash")
	{
		hashRouter.POST("/file", hashFile)
		hashRouter.POST("/[]string", h.hashStringMassive)
	}
}

// @Summary File
// @Tags hash
// @Description get file hash
// @ID hash-file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to hash"
// @Success 200 {object} hashOutput
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /hash/file [post]

type hashOutput struct {
	Hash string `json:"hash" form:"hash"`
}

func hashFile(c echo.Context) error {
	var file *multipart.FileHeader
	var err error
	if file, err = c.FormFile("file"); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse file from body")
	}
	fileReader, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file from body")
	}
	defer fileReader.Close()

	var hashString string
	if hashString, err = hash.CalculateFromFile(fileReader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to calculate hash for file")
	}
	return c.JSON(http.StatusOK, hashOutput{hashString})
}

// @Summary StringMassive
// @Tags hash
// @Description Get hash from a list of strings
// @ID hash-string-massive
// @Accept json
// @Produce json
// @Param body stringMassiveInput true "List of strings to hash"
// @Success 200 {object} hashOutput
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /hash/[]string [post]

type stringMassiveInput struct {
	Hashes []string `json:"hashes" form:"hashes" validate:"required"`
}

func (h *Handler) hashStringMassive(c echo.Context) error {
	var body stringMassiveInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse body")
	}

	if err := h.validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var hashString string
	var err error
	if hashString, err = hash.CalculateFromHashes(body.Hashes); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to calculate hash from hashes")
	}
	return c.JSON(http.StatusOK, hashOutput{hashString})
}
