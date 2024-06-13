package handler

type Deps struct {
	UserService  UserService
	FileService  FileService
	ChunkService ChunkService
}

type Handler struct {
	User  *UserHandler
	File  *FileHandler
	Chunk *ChunkHandler
	Hash  *HashHandler
	Error *ErrorHandler
}

func New(deps Deps) *Handler {
	return &Handler{
		User:  NewUserHandler(deps.UserService),
		File:  NewFileHandler(deps.FileService),
		Chunk: NewChunkHandler(deps.ChunkService, deps.FileService),
		Hash:  NewHashHandler(),
		Error: NewErrorHandler(),
	}
}
