package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/lapkomo2018/DiskordServer/internal/service"
	"github.com/lapkomo2018/DiskordServer/internal/storage"
	"github.com/lapkomo2018/DiskordServer/internal/storage/discord"
	"github.com/lapkomo2018/DiskordServer/internal/storage/gorm"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest/handler"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest/middleware"
	"github.com/lapkomo2018/DiskordServer/pkg/auth"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// @title Diskord API
// @version 1.0
// @description This is a sample swagger for Fiber

// @host localhost:3000
// @BasePath /api/

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	storages, err := storage.New(storage.Deps{
		GormDeps: gorm.Deps{
			DSN: os.Getenv("DB"),
		},
		DiscordDeps: discord.Deps{
			BotToken: os.Getenv("DISCORD_TOKEN"),
			Channel:  os.Getenv("DISCORD_CHANEL"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	tokenManager, err := auth.NewManager(os.Getenv("SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	services := service.New(service.Deps{
		UserStorage:    storages.Gorm.User,
		TokenManager:   tokenManager,
		AccessTokenTTL: time.Hour * 24 * 30,
	})

	middlewares := middleware.New(middleware.Deps{
		UserService:  services.User,
		FileService:  services.File,
		ChunkService: services.Chunk,
	})

	handlers := handler.New(handler.Deps{
		UserService:  services.User,
		FileService:  services.File,
		ChunkService: services.Chunk,
	})

	f := fiber.New(fiber.Config{BodyLimit: 1024 * 1024 * 25, ErrorHandler: handlers.Error.Handle})
	f.Use(logger.New())

	whiteListArray := strings.Split(os.Getenv("CORS_WHITELIST"), ",")
	for i, addr := range whiteListArray {
		whiteListArray[i] = strings.TrimSpace(addr)
	}
	f.Use(middleware.Cors(whiteListArray))

	f.Get("/swagger/*", swagger.HandlerDefault)

	api := f.Group("/api")

	hashGroup := api.Group("/hash")
	hashGroup.Post("/file", handlers.Hash.File)
	hashGroup.Post("/[]string", handlers.Hash.StringMassive)

	userGroup := api.Group("/user")
	userGroup.Post("/signup", handlers.User.Signup)
	userGroup.Post("/login", handlers.User.Login)
	userGroup.Get("/validate", middlewares.Auth.Require, handlers.User.Validate)
	userGroup.Get("/files", middlewares.Auth.Require, handlers.User.Files)

	filesGroup := api.Group("/files")
	filesGroup.Post("/upload", middlewares.Auth.Require, handlers.File.Upload)

	fileIdGroup := filesGroup.Group("/:fileId<min(1)>", middlewares.File.Require, skip.New(middlewares.Auth.Require, middlewares.File.IsPublic), skip.New(middlewares.File.OwnerCheck, middlewares.File.IsPublic))
	fileIdGroup.Get("/", handlers.File.Info)
	fileIdGroup.Get("/download", handlers.File.Download)
	fileIdGroup.Patch("/privacy", skip.New(middlewares.Auth.Require, middleware.IsKeyInLocals("user")), middlewares.File.OwnerCheck, handlers.File.ChangePrivacy)
	fileIdGroup.Delete("/", skip.New(middlewares.Auth.Require, middleware.IsKeyInLocals("user")), middlewares.File.OwnerCheck, handlers.File.Delete)

	chunksGroup := fileIdGroup.Group("/chunks")
	chunksGroup.Post("/upload", skip.New(middlewares.Auth.Require, middleware.IsKeyInLocals("user")), middlewares.File.OwnerCheck, handlers.Chunk.Upload)

	chunkIndexGroup := chunksGroup.Group("/:chunkIndex<min(0)>", middlewares.Chunk.Require)
	chunkIndexGroup.Get("/", handlers.Chunk.Info)
	chunkIndexGroup.Get("/download", handlers.Chunk.Download)

	log.Fatal(f.Listen(":3000"))
}
