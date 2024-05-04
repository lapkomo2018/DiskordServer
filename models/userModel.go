package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Files    []File `gorm:"foreignKey:UserID"`
}

type File struct {
	gorm.Model
	UserID uint    `gorm:"not null"`
	Name   string  `gorm:"not null"`
	Hash   string  `gorm:"not null"`
	Size   uint    `gorm:"not null"`
	Pieces []Piece `gorm:"foreignKey:FileID"`
	User   User    `gorm:"references:ID"`
}

type Piece struct {
	gorm.Model
	FileID    uint   `gorm:"not null"`
	MessageID string `gorm:"unique;not null"`
	Index     uint   `gorm:"not null"`
	File      File   `gorm:"references:ID"`
}
