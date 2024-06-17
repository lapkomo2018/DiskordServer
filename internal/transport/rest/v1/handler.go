package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/DiskordServer/internal/core"
	"io"
	"log"
	"mime/multipart"
)

type (
	UserService interface {
		Create(email, password string) (*core.User, error)
		Login(email, password string) (token string, err error)
		LoadFiles(user *core.User) error
		GetUserFromToken(jwtToken string) (*core.User, error)
	}

	FileService interface {
		First(file *core.File, cond ...interface{}) error
		Create(file *core.File) error
		Save(file *core.File) error
		Delete(file *core.File) error
		IsNeedChunk(file *core.File, chunkSize uint64) (bool, error)
		LoadChunks(file *core.File) error
	}

	ChunkService interface {
		First(chunk *core.Chunk, cond ...interface{}) error
		DownloadChunk(chunk *core.Chunk) (io.Reader, error)
		UploadChunk(chunk *core.Chunk, file *multipart.FileHeader) error
	}

	Validator interface {
		Struct(s interface{}) error
	}

	Handler struct {
		userService  UserService
		fileService  FileService
		chunkService ChunkService
		validator    Validator
	}
)

func New(userService UserService, fileService FileService, chunkService ChunkService, validator Validator) *Handler {
	return &Handler{
		userService:  userService,
		fileService:  fileService,
		chunkService: chunkService,
		validator:    validator,
	}
}

func (h *Handler) Init(api *echo.Group) {
	log.Printf("Initializing V1 api")
	h.initHashRouters(api)
	h.initUserRouters(api)
	h.initFilesRouters(api)
}
