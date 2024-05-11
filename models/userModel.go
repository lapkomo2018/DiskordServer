package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Files    []File `gorm:"foreignKey:UserId"`
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

type Chunk struct {
	gorm.Model
	FileId    uint   `gorm:"not null"`
	Index     uint   `gorm:"not null"`
	Hash      string `gorm:"not null"`
	Size      uint64 `gorm:"not null"`
	MessageId string `gorm:"unique;not null"`
	File      File   `gorm:"references:ID"`
}

func (p *Chunk) BeforeCreate(tx *gorm.DB) (err error) {
	if p.FileId == 0 {
		err = errors.New("file id cannot be zero")
	} else if p.MessageId == "" {
		err = errors.New("message id cannot be empty")
	} else if p.Hash == "" {
		err = errors.New("file hash cannot be empty")
	} else if p.Size == 0 {
		err = errors.New("file size cannot be zero")
	}
	return
}
