package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
)

type Deps struct {
	UserService  UserService
	FileService  FileService
	ChunkService ChunkService
}

type Middleware struct {
	Auth  *AuthMiddleware
	File  *FileMiddleware
	Chunk *ChunkMiddleware
}

func New(deps Deps) *Middleware {
	return &Middleware{
		Auth:  NewAuthMiddleware(deps.UserService),
		File:  NewFileMiddleware(deps.FileService),
		Chunk: NewChunkMiddleware(deps.ChunkService),
	}
}

func IsKeyInLocals(key string) func(c *fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		_, ok := c.Locals(key).(core.User)
		return ok
	}
}
