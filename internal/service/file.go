package service

import (
	"github.com/lapkomo2018/DiskordServer/internal/core"
)

type FileStorage interface {
	Create(file *core.File) error
	Save(file *core.File) error
	Delete(file *core.File) error
	Exists(id uint) error
	First(file *core.File, cond ...interface{}) error
	LoadChunks(file *core.File) error
}

type FileService struct {
	storage FileStorage
}

func NewFileService(s FileStorage) *FileService {
	return &FileService{storage: s}
}

func (s *FileService) Create(file *core.File) error {
	return s.storage.Create(file)
}

func (s *FileService) Save(file *core.File) error {
	return s.storage.Save(file)
}

func (s *FileService) Delete(file *core.File) error {
	return s.storage.Delete(file)
}

func (s *FileService) Exists(id uint) error {
	return s.storage.Exists(id)
}

func (s *FileService) First(file *core.File, cond ...interface{}) error {
	return s.storage.First(file, cond...)
}

func (s *FileService) IsNeedChunk(file *core.File, chunkSize uint64) (bool, error) {
	if err := s.LoadChunks(file); err != nil {
		return false, err
	}
	return file.IsNeedChunk(chunkSize), nil
}

func (s *FileService) LoadChunks(file *core.File) error {
	return s.storage.LoadChunks(file)
}
