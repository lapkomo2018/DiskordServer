package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/lapkomo2018/DiskordServer/controllers"
	"github.com/lapkomo2018/DiskordServer/handlers"
	"github.com/lapkomo2018/DiskordServer/initializers"
	"github.com/lapkomo2018/DiskordServer/middlewares"
	"log"
	"os"
	"strings"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.InitializeDiscordBot()
}

func main() {
	app := fiber.New(fiber.Config{BodyLimit: 1024 * 1024 * 25, ErrorHandler: handlers.ErrorHandler})
	app.Use(logger.New())

	corsWhiteList := os.Getenv("CORS_WHITELIST")
	whiteListArray := strings.Split(corsWhiteList, ",")
	for i, addr := range whiteListArray {
		whiteListArray[i] = strings.TrimSpace(addr)
	}
	app.Use(middlewares.CorsMiddleware(whiteListArray))

	app.Get("/monitor", monitor.New())

	api := app.Group("/api")

	hashes := api.Group("/hash")
	hashes.Post("/file", controllers.CalculateHashFromFile)
	hashes.Post("/[]string", controllers.CalculateHashFromHashes)

	user := api.Group("/user")
	user.Post("/signup", controllers.Signup)
	user.Post("/login", controllers.Login)
	user.Get("/validate", middlewares.RequireAuth, controllers.Validate)
	user.Get("/files", middlewares.RequireAuth, controllers.GetUserFiles)

	files := api.Group("/files")
	files.Post("/upload", middlewares.RequireAuth, controllers.UploadFile)

	fileId := files.Group("/:fileId<min(1)>", middlewares.RequireFile, skip.New(middlewares.RequireAuth, middlewares.FileIsPublic), skip.New(middlewares.FileOwnerCheck, middlewares.FileIsPublic))
	fileId.Get("/", controllers.GetFileInfo)
	fileId.Get("/download", controllers.DownloadFile)
	fileId.Patch("/privacy", skip.New(middlewares.RequireAuth, middlewares.IsKeyInLocals("user")), middlewares.FileOwnerCheck, controllers.ChangeFilePrivacy)
	fileId.Delete("/", skip.New(middlewares.RequireAuth, middlewares.IsKeyInLocals("user")), middlewares.FileOwnerCheck, controllers.DeleteFile)

	chunks := fileId.Group("/chunks")
	chunks.Post("/upload", skip.New(middlewares.RequireAuth, middlewares.IsKeyInLocals("user")), middlewares.FileOwnerCheck, controllers.UploadChunk)

	chunkIndex := chunks.Group("/:chunkIndex<min(0)>", middlewares.RequireChunk)
	chunkIndex.Get("/", controllers.GetChunkInfo)
	chunkIndex.Get("/download", controllers.DownloadChunk)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
