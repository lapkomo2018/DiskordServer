package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lapkomo2018/DiskordServer/controllers"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/middleware"
	"log"
	"os"
	"path/filepath"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	log.Println("Starting server...")
	app := fiber.New(fiber.Config{BodyLimit: 1024 * 1024 * 1024})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Static("/uploads", "./uploads")

	api := app.Group("/api")
	files := api.Group("/files")
	files.Post("/post", handleFilePost)

	user := api.Group("/user")
	user.Post("/signup", controllers.Signup)
	user.Post("/login", controllers.Login)
	user.Get("/validate", middleware.RequireAuth, controllers.Validate)

	userFile := user.Group("/file", middleware.RequireAuth)
	userFile.Get("/list", controllers.GetUserFilesList)
	userFile.Get("/download", Hello)
	userFile.Post("/send", Hello)

	userFilePiece := userFile.Group("/piece", middleware.RequireAuth)
	userFilePiece.Post("/send", Hello)
	userFilePiece.Get("/download", Hello)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func handleFilePost(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	path := filepath.Join("./uploads", file.Filename)
	if err := c.SaveFile(file, path); err != nil {
		return err
	}

	responseObject := fiber.Map{
		"fileUrl": "http://46.63.69.24:3000/uploads/" + file.Filename,
	}

	return c.JSON(responseObject)
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
