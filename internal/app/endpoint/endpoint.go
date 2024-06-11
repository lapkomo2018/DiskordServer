package endpoint

import (
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/chunk"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/chunks"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/file"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/files"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/hash"
	"github.com/lapkomo2018/DiskordServer/internal/app/endpoint/user"
	"github.com/lapkomo2018/DiskordServer/internal/model"
	"gorm.io/gorm"
)

type Endpoint struct {
	Chunk  *chunk.Chunk
	Chunks *chunks.Chunks
	File   *file.File
	Files  *files.Files
	Hash   *hash.Hash
	User   *user.User
}

func New(db *gorm.DB, discord model.DiscordService) *Endpoint {
	return &Endpoint{
		Chunk:  chunk.New(db, discord),
		Chunks: chunks.New(db, discord),
		File:   file.New(db),
		Files:  files.New(db),
		Hash:   hash.New(),
		User:   user.New(db),
	}
}
