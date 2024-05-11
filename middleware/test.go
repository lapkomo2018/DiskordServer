package middleware

import "github.com/gofiber/fiber/v2"

func TestParseBody(c *fiber.Ctx) error {
	var body struct {
		Hash  string
		Size  uint64
		Index uint64
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"body": body,
	})
}
