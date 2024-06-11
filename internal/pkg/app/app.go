package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	_ "github.com/lapkomo2018/DiskordServer/docs"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint"
	"github.com/lapkomo2018/DiskordServer/internal/app/handler"
	"github.com/lapkomo2018/DiskordServer/internal/app/middleware"
	"github.com/lapkomo2018/DiskordServer/internal/app/service"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

type App struct {
	f    *fiber.App
	port int
}

// @title Diskord API
// @version 1.0
// @description This is a sample swagger for Fiber

// @host localhost:3000
// @BasePath /api/

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func New(port int) (*App, error) {
	var err error
	app := &App{
		port: port,
	}

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(os.Getenv("DB")), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.User{}, &model.File{}, &model.Chunk{}); err != nil {
		return nil, err
	}

	s := service.New()
	if err := s.SetupDiscord(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_CHANEL")); err != nil {
		return nil, err
	}

	h := handler.New()
	m := middleware.New(db)
	e := endpoint.New(db, s.Discord)

	app.f = fiber.New(fiber.Config{BodyLimit: 1024 * 1024 * 25, ErrorHandler: h.Error.Handle})
	app.f.Use(fiberLogger.New())

	whiteListArray := strings.Split(os.Getenv("CORS_WHITELIST"), ",")
	for i, addr := range whiteListArray {
		whiteListArray[i] = strings.TrimSpace(addr)
	}
	app.f.Use(middleware.Cors(whiteListArray))

	app.f.Get("/swagger/*", swagger.HandlerDefault)

	app.f.Get("/monitor", monitor.New())

	api := app.f.Group("/api")

	hashGroup := api.Group("/hash")
	hashGroup.Post("/file", e.Hash.File)
	hashGroup.Post("/[]string", e.Hash.StringMassive)

	userGroup := api.Group("/user")
	userGroup.Post("/signup", e.User.Signup)
	userGroup.Post("/login", e.User.Login)
	userGroup.Get("/validate", m.Auth.Require, e.User.Validate)
	userGroup.Get("/files", m.Auth.Require, e.User.Files)

	filesGroup := api.Group("/files")
	filesGroup.Post("/upload", m.Auth.Require, e.Files.Upload)

	fileIdGroup := filesGroup.Group("/:fileId<min(1)>", m.File.Require, skip.New(m.Auth.Require, m.File.IsPublic), skip.New(m.File.OwnerCheck, m.File.IsPublic))
	fileIdGroup.Get("/", e.File.Info)
	fileIdGroup.Get("/download", e.File.Download)
	fileIdGroup.Patch("/privacy", skip.New(m.Auth.Require, middleware.IsKeyInLocals("user")), m.File.OwnerCheck, e.File.ChangePrivacy)
	fileIdGroup.Delete("/", skip.New(m.Auth.Require, middleware.IsKeyInLocals("user")), m.File.OwnerCheck, e.File.Delete)

	chunksGroup := fileIdGroup.Group("/chunks")
	chunksGroup.Post("/upload", skip.New(m.Auth.Require, middleware.IsKeyInLocals("user")), m.File.OwnerCheck, e.Chunks.Upload)

	chunkIndexGroup := chunksGroup.Group("/:chunkIndex<min(0)>", m.Chunk.Require)
	chunkIndexGroup.Get("/", e.Chunk.Info)
	chunkIndexGroup.Get("/download", e.Chunk.Download)

	return app, nil
}

func (a *App) Run() error {
	return a.f.Listen(fmt.Sprintf(":%d", a.port))
}
