package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/functions"
	"log"
	"mime/multipart"
)

func CalculateHashFromFile(c *fiber.Ctx) error {
	var file *multipart.FileHeader
	var err error
	if file, err = c.FormFile("file"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse file from body",
		})
	}
	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file from body",
		})
	}
	defer fileReader.Close()

	var hash string
	if hash, err = functions.CalculateHashFormFile(fileReader); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate hash for file",
		})
	}
	log.Printf("Returning hash from file: %s", hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": hash,
	})
}

func CalculateHashFromHashes(c *fiber.Ctx) error {
	var body struct {
		Hashes []string
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse body",
		})
	}

	var hash string
	var err error
	if hash, err = functions.CalculateHashFromHashes(body.Hashes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate hash from hashes",
		})
	}
	log.Printf("Returning hash from hashes: %s", hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"hash": hash,
	})
}
