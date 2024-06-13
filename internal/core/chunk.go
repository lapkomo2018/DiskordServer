package core

import (
	"errors"
	"gorm.io/gorm"
)

type Chunk struct {
	gorm.Model
	FileId    uint   `gorm:"not null"`
	Index     uint   `gorm:"not null"`
	Hash      string `gorm:"not null"`
	Size      uint64 `gorm:"not null"`
	MessageId string `gorm:"unique;not null"`
	File      File   `gorm:"references:ID"`
}

func (c *Chunk) BeforeCreate(tx *gorm.DB) (err error) {
	if c.FileId == 0 {
		err = errors.New("file id cannot be zero")
	} else if c.MessageId == "" {
		err = errors.New("message id cannot be empty")
	} else if c.Hash == "" {
		err = errors.New("file hash cannot be empty")
	} else if c.Size == 0 {
		err = errors.New("file size cannot be zero")
	}
	return
}
