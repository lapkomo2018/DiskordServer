package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Files    []File `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Email == "" {
		err = errors.New("email cannot be empty")
	} else if u.Password == "" {
		err = errors.New("password cannot be empty")
	}
	return
}

type File struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	Name      string  `gorm:"not null"`
	Hash      string  `gorm:"not null"`
	Size      uint64  `gorm:"not null"`
	IsPublic  bool    `gorm:"not null"`
	NumChunks uint    `gorm:"not null"`
	ChunkSize uint64  `gorm:"not null"`
	Chunks    []Chunk `gorm:"foreignKey:FileID"`
	User      User    `gorm:"references:ID"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.Size == 0 {
		err = errors.New("file size cannot be zero")
	} else if f.UserID == 0 {
		err = errors.New("user id cannot be zero")
	} else if f.Name == "" {
		err = errors.New("file name cannot be empty")
	} else if f.Hash == "" {
		err = errors.New("file hash cannot be empty")
	}
	return
}

type Chunk struct {
	gorm.Model
	FileID    uint   `gorm:"not null"`
	Index     uint   `gorm:"not null"`
	Hash      string `gorm:"not null"`
	Size      uint64 `gorm:"not null"`
	MessageID string `gorm:"unique;not null"`
	File      File   `gorm:"references:ID"`
}

func (p *Chunk) BeforeCreate(tx *gorm.DB) (err error) {
	if p.FileID == 0 {
		err = errors.New("file id cannot be zero")
	} else if p.MessageID == "" {
		err = errors.New("message id cannot be empty")
	} else if p.Hash == "" {
		err = errors.New("file hash cannot be empty")
	} else if p.Size == 0 {
		err = errors.New("file size cannot be zero")
	}
	return
}
