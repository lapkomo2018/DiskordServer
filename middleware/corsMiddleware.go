package middleware

import "github.com/gofiber/fiber/v2"

func CorsMiddleware(allowedOrigins []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		origin := c.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}
