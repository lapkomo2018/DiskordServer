package initializers

import "github.com/lapkomo2018/DiskordServer/models"

func SyncDatabase() {
	if err := DB.AutoMigrate(&models.User{}, &models.File{}, &models.Piece{}); err != nil {
		panic(err)
	}
}
