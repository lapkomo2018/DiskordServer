package discord

import "github.com/bwmarrin/discordgo"

type Deps struct {
	BotToken string
	Channel  string
}

type Storage struct {
	File *FileStorage
}

func New(deps Deps) (*Storage, error) {
	bot, err := discordgo.New("Bot " + deps.BotToken)
	if err != nil {
		return nil, err
	}

	return &Storage{
		File: NewFileStorage(bot, deps.Channel),
	}, nil
}
