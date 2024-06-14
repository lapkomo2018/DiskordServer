package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest/v1"
	"log"
	"net/http"
)

type ServicesV1 struct {
	UserService  v1.UserService
	FileService  v1.FileService
	ChunkService v1.ChunkService
}

type Deps struct {
	Port          int
	BodyLimit     int
	ServicesV1    ServicesV1
	CorsWhiteList []string
}

type Server struct {
	fiberApp   *fiber.App
	servicesV1 ServicesV1
	addr       string
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
		addr:       fmt.Sprintf(":%d", deps.Port),
		fiberApp:   f,
		servicesV1: deps.ServicesV1,
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
	handlerV1 := v1.New(s.servicesV1.UserService, s.servicesV1.FileService, s.servicesV1.ChunkService)
	api := s.fiberApp.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func (s *Server) Run() error {
	log.Println("Starting server")
	return s.fiberApp.Listen(s.addr)
}

func (s *Server) HttpServer() *http.Server {
	return &http.Server{
		Addr:    s.addr,
		Handler: s.HttpHandler(),
	}
}

func (s *Server) HttpHandler() http.Handler {
	return adaptor.FiberApp(s.fiberApp)
}
