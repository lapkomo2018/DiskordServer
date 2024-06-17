package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lapkomo2018/DiskordServer/internal/transport/rest/v1"
	"log"
	"strconv"
)

type ServicesV1 struct {
	v1.UserService
	v1.ChunkService
	v1.FileService
	v1.Validator
}

type Deps struct {
	Port          int
	BodyLimit     int
	ServicesV1    ServicesV1
	CorsWhiteList []string
}

type Server struct {
	echo       *echo.Echo
	servicesV1 ServicesV1
	addr       string
}

func New(deps Deps) *Server {
	log.Printf("Created server with port: %d", deps.Port)

	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler

	e.Use(middleware.BodyLimit(strconv.Itoa(deps.BodyLimit)))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())
	e.Use(Cors(deps.CorsWhiteList))

	e.Pre(middleware.RemoveTrailingSlash())

	return &Server{
		addr:       fmt.Sprintf(":%d", deps.Port),
		echo:       e,
		servicesV1: deps.ServicesV1,
	}
}

func (s *Server) Init() *Server {
	log.Println("Initializing server...")
	//s.fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	s.initApi()
	return s
}

func (s *Server) initApi() {
	log.Println("Initializing api...")
	handlerV1 := v1.New(s.servicesV1.UserService, s.servicesV1.FileService, s.servicesV1.ChunkService, s.servicesV1.Validator)
	api := s.echo.Group("/api")
	{
		handlerV1.Init(api)
	}
}

func (s *Server) Run() error {
	log.Println("Starting server")
	return s.echo.Start(s.addr)
}
