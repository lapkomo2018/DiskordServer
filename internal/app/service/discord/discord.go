package discord

import (
	"bytes"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
)

type Discord struct {
	channel string
	bot     *discordgo.Session
}

func New(token string, channel string) (*Discord, error) {
	var err error
	discord := &Discord{
		channel: channel,
	}
	discord.bot, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return discord, nil
}

func (d *Discord) UploadFile(fileName string, reader io.Reader) (string, error) {
	message, err := d.bot.ChannelFileSend(d.channel, fileName, reader)
	if err != nil {
		return "", err
	}
	return message.ID, nil
}

func (d *Discord) DownloadFileFromMessage(messageId string) (io.Reader, error) {
	message, err := d.bot.ChannelMessage(d.channel, messageId)
	if err != nil {
		return nil, err
	}

	for _, attachment := range message.Attachments {
		response, err := http.Get(attachment.URL)
		if err != nil {
			return nil, err
		}

		var buffer bytes.Buffer

		if _, err := io.Copy(&buffer, response.Body); err != nil {
			response.Body.Close()
			return nil, err
		}
		response.Body.Close()

		bufferReader := bytes.NewReader(buffer.Bytes())

		return bufferReader, nil
	}

	return nil, errors.New("no attachment found")
}

func (d *Discord) DeleteMessage(messageId string) error {
	if err := d.bot.ChannelMessageDelete(d.channel, messageId); err != nil {
		return err
	}
	return nil
}
