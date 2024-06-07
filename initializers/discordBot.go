package initializers

import (
	"bytes"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"os"
)

var DiscordBot *discordgo.Session

func InitializeDiscordBot() {
	var err error
	DiscordBot, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}
}

func DownloadFileFromMessage(messageID string) (io.Reader, error) {
	message, err := DiscordBot.ChannelMessage(os.Getenv("DISCORD_CHANEL"), messageID)
	if err != nil {
		return nil, err
	}

	for _, attachment := range message.Attachments {
		response, err := http.Get(attachment.URL)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		var buffer bytes.Buffer

		if _, err := io.Copy(&buffer, response.Body); err != nil {
			return nil, err
		}

		bufferReader := bytes.NewReader(buffer.Bytes())

		return bufferReader, nil
	}

	return nil, errors.New("No attachment found")
}

func DeleteMessage(messageID string) error {
	if err := DiscordBot.ChannelMessageDelete(os.Getenv("DISCORD_CHANEL"), messageID); err != nil {
		return err
	}
	return nil
}
