package gv

import (
	"errors"
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo

	cfg ServerConfig
}

func (s *Server) Build() error {
	result := api.Build(s.cfg.EsBuildOptions)
	if len(result.Errors) > 0 {
		return errors.New(result.Errors[0].Text)
	}

	return nil
}

func (s *Server) Server() *echo.Echo {
	return s.e
}

func (s *Server) Start() error {
	err := s.Build()
	if err != nil {
		return err
	}

	fmt.Println("\n" + color.GreenString("➜") + " Local: " + color.MagentaString(fmt.Sprintf("http://%s:%d", s.cfg.Host, s.cfg.Port)) + "\n")

	return s.e.Start(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
}

type ServerConfig struct {
	Host string
	Port int

	EsBuildOptions api.BuildOptions
}

func NewServer(cfg ServerConfig) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		e:   e,
		cfg: cfg,
	}
}
