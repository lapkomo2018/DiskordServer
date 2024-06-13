package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"strings"
)

type UserService interface {
	Auth(jwtToken string) (*core.User, error)
}

type AuthMiddleware struct {
	userService UserService
}

func NewAuthMiddleware(userService UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (a *AuthMiddleware) Require(c *fiber.Ctx) error {
	// Get cookie
	tokenString, err := extractTokenFromCookies(c)
	if err != nil {
		tokenString, err = extractTokenFromHeaders(c)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Failed to get token")
		}
	}

	user, err := a.userService.Auth(tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Locals("user", *user)

	return c.Next()
}

func extractTokenFromCookies(c *fiber.Ctx) (string, error) {
	token := c.Cookies("Authorization")
	if len(token) == 0 {
		return "", errors.New("no token in cookies")
	}
	return token, nil
}

func extractTokenFromBody(c *fiber.Ctx) (string, error) {
	var body struct {
		Authorization string
	}
	if err := c.BodyParser(&body); err != nil {
		return "", errors.New("failed to parse body")
	}
	if len(body.Authorization) == 0 {
		return "", errors.New("no token in body")
	}
	return body.Authorization, nil
}

func extractTokenFromHeaders(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
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
