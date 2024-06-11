package hash

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/pkg/hash"
	"mime/multipart"
)

type Hash struct {
}

func New() *Hash {
	return &Hash{}
}

type hashOutput struct {
	Hash string `json:"hash"`
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

func (h *Hash) File(c *fiber.Ctx) error {
	var file *multipart.FileHeader
	var err error
	if file, err = c.FormFile("file"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse file from body")
	}
	fileReader, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file from body")
	}
	defer fileReader.Close()

	var hashString string
	if hashString, err = hash.CalculateFromFile(fileReader); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate hash for file")
	}
	return c.Status(fiber.StatusOK).JSON(hashOutput{hashString})
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
	Hashes []string `json:"hashes"`
}

func (h *Hash) StringMassive(c *fiber.Ctx) error {
	var body stringMassiveInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	var hashString string
	var err error
	if hashString, err = hash.CalculateFromHashes(body.Hashes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate hash from hashes")
	}
	return c.Status(fiber.StatusOK).JSON(hashOutput{hashString})
}
