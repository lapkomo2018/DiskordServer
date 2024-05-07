package initializers

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"os"
)

var DiscordBot *discordgo.Session

const ChannelID string = "1237128843039604918"

func InitializeDiscordBot() {
	var err error
	DiscordBot, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}
}

func DownloadFilesFromMessage(channelID string, messageID string) (io.Reader, error) {
	message, err := DiscordBot.ChannelMessage(channelID, messageID)
	if err != nil {
		return nil, err
	}

	for _, attachment := range message.Attachments {
		response, err := http.Get(attachment.URL)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		return response.Body, nil
	}

	return nil, errors.New("No attachment found")
}
