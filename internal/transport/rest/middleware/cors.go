package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Cors(allowedOrigins []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		origin := c.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if strings.HasPrefix(origin, "http://"+allowedOrigin) || strings.HasPrefix(origin, "https://"+allowedOrigin) {
				c.Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}
