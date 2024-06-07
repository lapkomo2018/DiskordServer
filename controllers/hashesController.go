package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/functions"
	"mime/multipart"
)

func CalculateHashFromFile(c *fiber.Ctx) error {
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

	var hash string
	if hash, err = functions.CalculateHashFormFile(fileReader); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate hash for file")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": hash,
	})
}

func CalculateHashFromHashes(c *fiber.Ctx) error {
	var body struct {
		Hashes []string
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse body")
	}

	var hash string
	var err error
	if hash, err = functions.CalculateHashFromHashes(body.Hashes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate hash from hashes")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": hash,
	})
}
