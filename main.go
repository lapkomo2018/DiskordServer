package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lapkomo2018/DiskordServer/controllers"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/middleware"
	"log"
	"os"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.InitializeDiscordBot()
}

func main() {
	app := fiber.New(fiber.Config{BodyLimit: 1024 * 1024 * 25})
	app.Use(middleware.CorsMiddleware([]string{"http://46.63.69.24:5173", "http://localhost:5173"}))

	api := app.Group("/api")

	user := api.Group("/user")
	user.Post("/signup", controllers.Signup)
	user.Post("/login", controllers.Login)
	user.Get("/validate", middleware.RequireAuth, controllers.Validate)

	userFile := user.Group("/file", middleware.RequireAuth)
	userFile.Get("/list", controllers.GetUserFilesList)
	userFile.Post("/upload", controllers.UploadFile)
	userFile.Get("/download", middleware.FileAccessCheck, controllers.DownloadFile)

	userFileChunk := userFile.Group("/chunk", middleware.RequireAuth, middleware.FileAccessCheck)
	userFileChunk.Post("/upload", controllers.UploadChunk)
	userFileChunk.Get("/download", controllers.DownloadChunk)

	hashes := api.Group("/hashes")
	hashes.Post("/file", controllers.CalculateHashFromFile)
	hashes.Post("/hashes", controllers.CalculateHashFromHashes)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
