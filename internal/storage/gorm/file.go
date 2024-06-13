package gorm

import (
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"gorm.io/gorm"
)

type FileStorage struct {
	db *gorm.DB
}

func NewFileStorage(db *gorm.DB) *FileStorage {
	return &FileStorage{
		db: db,
	}
}

// TODO: FIX ERROR
func (us *FileStorage) First(file *core.File, cond ...interface{}) error {
	return us.db.First(file, cond...).Error
}

func (us *FileStorage) FindAll(dest interface{}, conds ...interface{}) error {
	return us.db.Find(dest, conds...).Error
}

func (us *FileStorage) Exists(id uint) error {
	return us.db.First(core.File{}, id).Error
}

func (us *FileStorage) Create(file *core.File) error {
	return us.db.Create(file).Error
}

func (us *FileStorage) Save(file *core.File) error {
	return us.db.Save(file).Error
}

func (us *FileStorage) Delete(file *core.File) error {
	return us.db.Delete(file).Error
}

func (us *FileStorage) LoadChunks(file *core.File) error {
	return us.db.Preload("Chunks").First(file).Error
}
