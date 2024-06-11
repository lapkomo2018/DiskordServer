package service

import (
	"github.com/lapkomo2018/DiskordServer/internal/app/service/discord"
)

type Service struct {
	Discord *discord.Discord
}

func New() *Service {
	return &Service{}
}

func (s *Service) SetupDiscord(discordToken string, discordChannel string) error {
	d, err := discord.New(discordToken, discordChannel)
	if err != nil {
		return err
	}

	s.Discord = d
	return nil
}
