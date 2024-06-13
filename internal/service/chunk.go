package service

import (
	"bytes"
	"errors"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"github.com/lapkomo2018/DiskordServer/pkg/hash"
	"io"
	"mime/multipart"
)

type ChunkStorage interface {
	Create(chunk *core.Chunk) error
	Exists(id uint) error
	First(chunk *core.Chunk, cond ...interface{}) error
}

type DiscordFileStorage interface {
	UploadFile(fileName string, reader io.Reader) (string, error)
	DownloadFileFromMessage(messageId string) (io.Reader, error)
}

type ChunkService struct {
	dbStorage   ChunkStorage
	fileStorage DiscordFileStorage
}

func NewChunkService(s ChunkStorage) *ChunkService {
	return &ChunkService{dbStorage: s}
}

func (s *ChunkService) Exists(id uint) error {
	return s.dbStorage.Exists(id)
}

func (s *ChunkService) First(chunk *core.Chunk, cond ...interface{}) error {
	return s.dbStorage.First(chunk, cond...)
}

func (s *ChunkService) DownloadChunk(chunk core.Chunk) (io.Reader, error) {
	return s.fileStorage.DownloadFileFromMessage(chunk.MessageId)
}

func (s *ChunkService) UploadChunk(chunk *core.Chunk, file *multipart.FileHeader) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return err
	}

	hashString, err := hash.CalculateFromFile(bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}
	if hashString != chunk.Hash {
		return errors.New("hash is invalid")
	}

	messageId, err := s.fileStorage.UploadFile(file.Filename, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	chunk.MessageId = messageId

	if err := s.dbStorage.Create(chunk); err != nil {
		return err
	}

	return nil
}
