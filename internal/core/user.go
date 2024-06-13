package core

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

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("user_id = ?", u.ID).Delete(&File{})
	return
}
