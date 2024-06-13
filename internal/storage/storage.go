package storage

import (
	"github.com/lapkomo2018/DiskordServer/internal/storage/discord"
	"github.com/lapkomo2018/DiskordServer/internal/storage/gorm"
)

type Deps struct {
	GormDeps    gorm.Deps
	DiscordDeps discord.Deps
}

type Storage struct {
	Gorm    *gorm.Storage
	Discord *discord.Storage
}

func New(deps Deps) (*Storage, error) {
	discordStorage, err := discord.New(deps.DiscordDeps)
	if err != nil {
		return nil, err
	}

	gormStorage, err := gorm.New(deps.GormDeps)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Discord: discordStorage,
		Gorm:    gormStorage,
	}, nil
}
