package gorm

import (
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Deps struct {
	Dsn string
}

type Storage struct {
	User  *UserStorage
	File  *FileStorage
	Chunk *ChunkStorage
}

func New(deps Deps) (*Storage, error) {
	log.Printf("Connecting db...")
	db, err := gorm.Open(postgres.Open(deps.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&core.User{}, &core.File{}, &core.Chunk{}); err != nil {
		return nil, err
	}

	return &Storage{
		User:  NewUserStorage(db),
		File:  NewFileStorage(db),
		Chunk: NewChunkStorage(db),
	}, nil
}
