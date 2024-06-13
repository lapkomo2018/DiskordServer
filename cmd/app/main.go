package main

import (
	"github.com/joho/godotenv"
	"github.com/lapkomo2018/DiskordServer/internal/service"
	"github.com/lapkomo2018/DiskordServer/internal/storage"
	"github.com/lapkomo2018/DiskordServer/internal/storage/discord"
	"github.com/lapkomo2018/DiskordServer/internal/storage/gorm"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest"
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
			Dsn: os.Getenv("DB"),
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
		UserStorage:        storages.Gorm.User,
		FileStorage:        storages.Gorm.File,
		ChunkStorage:       storages.Gorm.Chunk,
		DiscordFileStorage: storages.Discord.File,
		TokenManager:       tokenManager,
		AccessTokenTTL:     time.Hour * 24 * 30,
	})

	corsWhiteList := strings.Split(os.Getenv("CORS_WHITELIST"), ",")
	for i, addr := range corsWhiteList {
		corsWhiteList[i] = strings.TrimSpace(addr)
	}

	httpServer := rest.New(rest.Deps{
		Services:      services,
		BodyLimit:     1024 * 1024 * 25,
		Port:          3000,
		CorsWhiteList: corsWhiteList,
	}).Init()

	log.Fatal(httpServer.Run())
}
