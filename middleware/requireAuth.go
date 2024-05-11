package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/models"
	"os"
	"strings"
	"time"
)

func RequireAuth(c *fiber.Ctx) error {
	// Get cookie
	tokenString, err := extractTokenFromCookies(c)
	if err != nil {
		tokenString, err = extractTokenFromHeaders(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to get token",
			})
		}
	}

	// Decode/validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to decode token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is expired",
			})
		}
		//find user with token
		var user models.User
		if err := initializers.DB.First(&user, claims["sub"]).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		//attach to req
		c.Locals("user", user)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to parse token",
		})
	}
	//continue
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
