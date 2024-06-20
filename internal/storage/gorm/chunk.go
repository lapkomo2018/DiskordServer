package gorm

import (
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"gorm.io/gorm"
)

type ChunkStorage struct {
	db *gorm.DB
}

func NewChunkStorage(db *gorm.DB) *ChunkStorage {
	return &ChunkStorage{
		db: db,
	}
}

func (us *ChunkStorage) First(chunk *core.Chunk, cond ...interface{}) error {
	return us.db.Where(chunk).First(chunk, cond...).Error
}

func (us *ChunkStorage) FindAll(dest interface{}, conds ...interface{}) error {
	return us.db.Find(dest, conds...).Error
}

func (us *ChunkStorage) Exists(id uint) error {
	return us.db.First(&core.Chunk{}, id).Error
}

func (us *ChunkStorage) Create(chunk *core.Chunk) error {
	return us.db.Create(chunk).Error
}
