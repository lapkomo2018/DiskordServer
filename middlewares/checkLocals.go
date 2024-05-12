package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/models"
)

func IsKeyInLocals(key string) func(c *fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		_, ok := c.Locals(key).(models.User)
		return ok
	}
}
