package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	api := e.Group("/api")
	{
		private := api.Group("", privateMiddle)
		{
			private.GET("/private", privateHandler)
		}

		public := api.Group("", publicMiddle)
		{
			public.GET("/public", publicHandler)
		}
	}

	e.Logger.Fatal(e.Start(":3000"))

}

func privateMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("privateMiddle")
		return next(c)
	}
}

func publicMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("publicMiddle")
		return next(c)
	}
}

func privateHandler(c echo.Context) error {
	return c.String(http.StatusOK, "private")
}
func publicHandler(c echo.Context) error {
	var fileReader io.Reader
	return c.Stream(http.StatusOK, "text/plain", fileReader)
	return c.String(http.StatusOK, "public")
}
