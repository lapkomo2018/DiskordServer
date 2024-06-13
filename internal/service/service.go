package service

import "time"

type Deps struct {
	UserStorage  UserStorage
	FileStorage  FileStorage
	ChunkStorage ChunkStorage

	TokenManager   TokenManager
	AccessTokenTTL time.Duration
}

type Service struct {
	User  *UserService
	File  *FileService
	Chunk *ChunkService
}

func New(deps Deps) *Service {
	return &Service{
		User:  NewUserService(deps.UserStorage, deps.TokenManager, deps.AccessTokenTTL),
		File:  NewFileService(deps.FileStorage),
		Chunk: NewChunkService(deps.ChunkStorage),
	}
}
