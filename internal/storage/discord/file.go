package discord

import (
	"bytes"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
)

type FileStorage struct {
	channel string
	bot     *discordgo.Session
}

func NewFileStorage(bot *discordgo.Session, channel string) *FileStorage {
	return &FileStorage{
		channel: channel,
		bot:     bot,
	}
}

func (d *FileStorage) UploadFile(fileName string, reader io.Reader) (string, error) {
	message, err := d.bot.ChannelFileSend(d.channel, fileName, reader)
	if err != nil {
		return "", err
	}
	return message.ID, nil
}

func (d *FileStorage) DownloadFileFromMessage(messageId string) (io.Reader, error) {
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
		_, err = io.Copy(&buffer, response.Body)
		response.Body.Close()
		if err != nil {
			return nil, err
		}

		bufferReader := bytes.NewReader(buffer.Bytes())

		return bufferReader, nil
	}

	return nil, errors.New("no attachment found")
}
