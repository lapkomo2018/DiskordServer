package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationCookie = "Authorization"

	FileIdParams     = "fileId"
	ChunkIndexParams = "chunkIndex"

	CookieExpireTTL = time.Hour * 24 * 30

	UserLocals  = "user"
	FileLocals  = "file"
	ChunkLocals = "chunk"

	BearerSchema = "Bearer "
)

func (h *Handler) userIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := h.setUserFromRequest(c)
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (h *Handler) setUserFromRequest(c echo.Context) error {
	token, err := extractAuthTokenFromCookies(c)
	if err != nil {
		token, err = extractAuthTokenFromHeaders(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Failed to get token")
		}
	}

	user, err := h.userService.GetUserFromToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	c.Set(UserLocals, *user)
	return nil
}

func extractAuthTokenFromCookies(c echo.Context) (string, error) {
	tokenCookie, err := c.Cookie(AuthorizationCookie)
	if err != nil {
		return "", err
	}
	return tokenCookie.Value, nil
}

func extractAuthTokenFromHeaders(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get(AuthorizationHeader)
	if authHeader == "" {
		return "", errors.New("no token in header")
	}
	if !strings.HasPrefix(authHeader, BearerSchema) {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	token := strings.TrimPrefix(authHeader, BearerSchema)
	return token, nil
}

func (h *Handler) fileIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := h.setFileFromRequest(c)
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (h *Handler) setFileFromRequest(c echo.Context) error {
	fileIdString := c.Param(FileIdParams)
	if fileIdString == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file ID")
	}

	fileId, err := strconv.Atoi(fileIdString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file ID")
	}

	file := &core.File{
		ID: uint(fileId),
	}
	if err := h.fileService.First(file); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}

	c.Set(FileLocals, *file)
	return nil
}

func (h *Handler) fileAccessCheck(allowPublic bool) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file, ok := c.Get(FileLocals).(core.File)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse file")
			}

			if allowPublic && file.IsPublic {
				return next(c)
			}

			if !isInLocals(UserLocals, c) {
				if err := h.setUserFromRequest(c); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}

			user, ok := c.Get(UserLocals).(core.User)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse user")
			}

			if file.UserId != user.ID {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			}

			return next(c)
		}
	}
}

func (h *Handler) chunkIdentify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := h.setChunkFromRequest(c)
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (h *Handler) setChunkFromRequest(c echo.Context) error {
	file, ok := c.Get(FileLocals).(core.File)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse BD file")
	}

	chunkIndex, err := strconv.Atoi(c.Param(ChunkIndexParams))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chunk index")
	}

	chunk := &core.Chunk{
		FileId: file.ID,
		Index:  uint(chunkIndex),
	}
	if err := h.chunkService.First(chunk); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Chunk not found")
	}

	c.Set(ChunkLocals, *chunk)
	return nil
}

func isInLocals(key string, c echo.Context) bool {
	value := c.Get(key)
	return value != nil
}
