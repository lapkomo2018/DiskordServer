package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/skip"
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
	app.Use(logger.New())

	app.Get("/monitor", monitor.New())

	api := app.Group("/api")

	hashes := api.Group("/hash")
	hashes.Post("/file", controllers.CalculateHashFromFile)
	hashes.Post("/[]string", controllers.CalculateHashFromHashes)

	user := api.Group("/user")
	user.Post("/signup", controllers.Signup)
	user.Post("/login", controllers.Login)
	user.Get("/validate", middleware.RequireAuth, controllers.Validate)
	user.Get("/files", middleware.RequireAuth, controllers.GetUserFiles)

	file := api.Group("/file")
	file.Post("/upload", middleware.RequireAuth, controllers.UploadFile)

	fileId := file.Group("/:fileId<min(1)>", middleware.RequireFile, skip.New(middleware.RequireAuth, middleware.FileIsPublic), skip.New(middleware.FileOwnerCheck, middleware.FileIsPublic))
	fileId.Get("/info", controllers.GetFileInfo)
	fileId.Get("/download", controllers.DownloadFile)

	chunk := fileId.Group("/chunk")
	chunk.Post("/upload", skip.New(middleware.RequireAuth, middleware.IsKeyInLocals("user")), middleware.FileOwnerCheck, controllers.UploadChunk)

	chunkIndex := chunk.Group("/:chunkIndex<min(0)>", middleware.RequireChunk)
	chunkIndex.Get("/info", controllers.GetChunkInfo)
	chunkIndex.Get("/download", controllers.DownloadChunk)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
