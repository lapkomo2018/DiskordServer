package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Cors(allowedOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			origin := c.Request().Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if strings.HasPrefix(origin, "http://"+allowedOrigin) || strings.HasPrefix(origin, "https://"+allowedOrigin) {
					c.Response().Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
