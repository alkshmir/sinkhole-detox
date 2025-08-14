package presentation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alkshmir/sinkhole-detox.git/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	Port uint
}

type Server struct {
	e         *echo.Echo
	config    ServerConfig
	generator *domain.HostsGenerator
}

func NewServer(b []domain.Blocker, conf ServerConfig) *Server {
	generator := domain.NewHostsGenerator(b)
	e := echo.New()
	s := &Server{
		e:         e,
		config:    conf,
		generator: generator,
	}

	e.Use(middleware.Logger())

	e.GET("/", s.genHosts)

	return s
}

func (s *Server) Start() {
	port := s.config.Port
	if port == 0 {
		port = 8080
	}
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%d", port)))
}

func (s *Server) genHosts(c echo.Context) error {
	t := time.Now()
	entries := s.generator.Gen(t)
	var response string
	for _, entry := range entries {
		response += entry.String() + "\n"
	}
	return c.String(http.StatusOK, response)
}
