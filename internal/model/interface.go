package model

import "io"

type DiscordService interface {
	UploadFile(fileName string, reader io.Reader) (string, error)
	DownloadFileFromMessage(messageID string) (io.Reader, error)
	DeleteMessage(messageID string) error
}
