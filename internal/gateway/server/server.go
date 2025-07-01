package server

import (
	"fmt"
	"net/http"

	"ebus/internal/gateway/api"

	"go.uber.org/zap"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct {
	s      *http.Server
	config Config
	logger *zap.Logger
}

func NewServer(logger *zap.Logger, config Config, handler *api.API) (*Server, error) {
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: handler.GetHandler(),
	}

	return &Server{
		s:      s,
		config: config,
		logger: logger,
	}, nil
}

func (s *Server) Run() error {
	s.logger.Info(
		"Running HTTP Server",
		zap.String("host", s.config.Host),
		zap.Int("port", s.config.Port),
	)
	return s.s.ListenAndServe()
}
