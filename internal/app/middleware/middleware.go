package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/app/middleware/auth"
	"github.com/lapkomo2018/DiskordServer/internal/app/middleware/chunk"
	"github.com/lapkomo2018/DiskordServer/internal/app/middleware/file"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
	"strings"
)

type Middleware struct {
	Auth  *auth.Auth
	Chunk *chunk.Chunk
	File  *file.File
}

func New(db *gorm.DB) *Middleware {
	return &Middleware{
		Auth:  auth.New(db),
		Chunk: chunk.New(db),
		File:  file.New(db),
	}
}

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

func IsKeyInLocals(key string) func(c *fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		_, ok := c.Locals(key).(model.User)
		return ok
	}
}
