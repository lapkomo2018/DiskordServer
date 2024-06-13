package service

import (
	"log"
	"time"
)

type Deps struct {
	UserStorage  UserStorage
	FileStorage  FileStorage
	ChunkStorage ChunkStorage

	DiscordFileStorage DiscordFileStorage

	TokenManager   TokenManager
	AccessTokenTTL time.Duration
}

type Service struct {
	User  *UserService
	File  *FileService
	Chunk *ChunkService
}

func New(deps Deps) *Service {
	log.Printf("Creating services...")
	return &Service{
		User:  NewUserService(deps.UserStorage, deps.TokenManager, deps.AccessTokenTTL),
		File:  NewFileService(deps.FileStorage),
		Chunk: NewChunkService(deps.ChunkStorage, deps.DiscordFileStorage),
	}
}
