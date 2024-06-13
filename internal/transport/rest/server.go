package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lapkomo2018/DiskordServer/internal/service"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest/v1"
	"log"
)

type Deps struct {
	Port          int
	BodyLimit     int
	Services      *service.Service
	CorsWhiteList []string
}

type Server struct {
	fiberApp *fiber.App
	services *service.Service
	addr     string
}

func New(deps Deps) *Server {
	log.Printf("Created server with port: %d", deps.Port)
	f := fiber.New(fiber.Config{
		BodyLimit:    deps.BodyLimit,
		ErrorHandler: ErrorHandler,
	})

	f.Use(logger.New())
	f.Use(Cors(deps.CorsWhiteList))

	return &Server{
		addr:     fmt.Sprintf(":%d", deps.Port),
		fiberApp: f,
		services: deps.Services,
	}
}

func (s *Server) Init() *Server {
	log.Println("Initializing server...")
	s.fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	s.initApi()
	return s
}

func (s *Server) initApi() {
	log.Println("Initializing api...")
	handlerV1 := v1.New(v1.Deps{
		UserService:  s.services.User,
		FileService:  s.services.File,
		ChunkService: s.services.Chunk,
	})
	api := s.fiberApp.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func (s *Server) Run() error {
	log.Println("Starting server")
	return s.fiberApp.Listen(s.addr)
}
