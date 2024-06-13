package v1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"strconv"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationCookie = "Authorization"

	FileIdParams     = "fileId"
	ChunkIndexParams = "chunkIndex"

	UserLocals  = "user"
	FileLocals  = "file"
	ChunkLocals = "chunk"
)

func (h *Handler) userIdentify(c *fiber.Ctx) error {
	token, err := extractTokenFromCookies(c)
	if err != nil {
		token, err = extractTokenFromHeaders(c)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Failed to get token")
		}
	}

	user, err := h.userService.GetUserFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Locals(UserLocals, *user)

	return c.Next()
}

func extractTokenFromCookies(c *fiber.Ctx) (string, error) {
	token := c.Cookies(AuthorizationCookie)
	if len(token) == 0 {
		return "", errors.New("no token in cookies")
	}
	return token, nil
}

func extractTokenFromHeaders(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(AuthorizationHeader)
	if len(authHeader) == 0 {
		return "", errors.New("no token in header")
	}

	const bearerSchema = "Bearer "
	if !strings.HasPrefix(authHeader, bearerSchema) {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	token := strings.TrimPrefix(authHeader, bearerSchema)
	return token, nil
}

func (h *Handler) setFileFromRequest(c *fiber.Ctx) error {
	fileId := c.Params(FileIdParams)
	if fileId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file ID")
	}

	file := &core.File{}
	if err := h.fileService.First(file, fileId); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	c.Locals(FileLocals, *file)
	return c.Next()
}

func fileOwnerCheck(c *fiber.Ctx) error {
	user, ok := c.Locals(UserLocals).(core.User)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user")
	}

	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse file")
	}

	if file.UserId != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}
	return c.Next()
}

func fileIsPublic(c *fiber.Ctx) bool {
	file, ok := c.Locals(FileLocals).(core.File)
	if !ok {
		return false
	}
	return file.IsPublic
}

func (h *Handler) setChunkFromRequest(ctx *fiber.Ctx) error {
	var err error
	file, ok := ctx.Locals(FileLocals).(core.File)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse BD file")
	}

	var chunkIndex int
	chunkIndex, err = strconv.Atoi(ctx.Params(ChunkIndexParams))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid chunk index")
	}

	chunk := &core.Chunk{
		FileId: file.ID,
		Index:  uint(chunkIndex),
	}
	if err := h.chunkService.First(chunk); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chunk not found")
	}

	ctx.Locals(ChunkLocals, *chunk)
	return ctx.Next()
}

func isInLocals(key string) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		value := c.Locals(key)
		return value != nil
	}
}
