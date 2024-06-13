package core

import (
	"errors"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	UserId    uint    `gorm:"not null"`
	Name      string  `gorm:"not null"`
	Hash      string  `gorm:"not null"`
	Size      uint64  `gorm:"not null"`
	IsPublic  bool    `gorm:"not null"`
	NumChunks uint    `gorm:"not null"`
	ChunkSize uint64  `gorm:"not null"`
	Chunks    []Chunk `gorm:"foreignKey:FileId"`
	User      User    `gorm:"references:ID"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.Size == 0 {
		err = errors.New("file size cannot be zero")
	} else if f.UserId == 0 {
		err = errors.New("user id cannot be zero")
	} else if f.Name == "" {
		err = errors.New("file name cannot be empty")
	} else if f.Hash == "" {
		err = errors.New("file hash cannot be empty")
	}
	return
}

func (f *File) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("file_id = ?", f.ID).Delete(&Chunk{})
	return
}

func (f *File) Validate() bool {
	if f.NumChunks != uint(len(f.Chunks)) {
		return false
	}

	var totalSize uint64
	for _, chunk := range f.Chunks {
		totalSize += chunk.Size
	}
	if totalSize != f.Size {
		return false
	}

	return true
}

func (f *File) IsNeedChunk(chunkSize uint64) bool {
	var currentSize uint64
	for _, chunk := range f.Chunks {
		currentSize += chunk.Size
	}

	remainingSize := f.Size - currentSize

	return remainingSize >= chunkSize
}
