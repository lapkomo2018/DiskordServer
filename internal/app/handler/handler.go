package handler

import "github.com/lapkomo2018/DiskordServer/internal/app/handler/error"

type Handler struct {
	Error *error.Error
}

func New() *Handler {
	return &Handler{
		Error: error.New(),
	}
}
