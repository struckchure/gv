package gv

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo

	cfg ServerConfig
}

func (s *Server) Watch() error {
	ctx, esErr := api.Context(s.cfg.EsBuildOptions)
	if esErr != nil {
		return esErr
	}

	if err := ctx.Watch(api.WatchOptions{}); err != nil {
		return err
	}

	err := s.HandleHMR()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Build() error {
	ctx, esErr := api.Context(s.cfg.EsBuildOptions)
	if esErr != nil {
		return esErr
	}
	res := ctx.Rebuild()
	if len(res.Errors) > 0 {
		return errors.New(res.Errors[0].Text)
	}
	return nil
}

func (s *Server) Server() *echo.Echo {
	return s.e
}

var HmrClientBroadcast = make(chan HmrResult)

func (s *Server) HandleHMR() error {
	color.Magenta("[HMR] Wating for client to connect ...")

	var upgrader = websocket.Upgrader{}

	handler := func(c echo.Context) error {
		// Upgrade connection to WebSocket
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		color.Green("[HMR] Client connected")

		for {
			for payload := range HmrClientBroadcast {
				if err := ws.WriteJSON(payload); err != nil {
					c.Logger().Error("WebSocket broadcast error:", err)
					continue
				}
			}
		}
	}

	s.e.GET("/__hmr__", handler)

	return nil
}

func (s *Server) Start() error {
	mode := os.Getenv("GV_MODE")

	switch mode {
	case "dev":
		s.Watch()
	case "build":
		if err := s.Build(); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	fmt.Println("\n" + color.GreenString("➜") + " Local: " + color.MagentaString(fmt.Sprintf("http://%s:%d", s.cfg.Host, s.cfg.Port)) + "\n")

	return s.e.Start(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
}

type ServerConfig struct {
	Host string
	Port int

	EsBuildOptions    api.BuildOptions
	WatchPath         *string
	WatchExcludePaths *[]string
}

func NewServer(cfg ServerConfig) *Server {
	if os.Getenv("GV_MODE") == "" {
		os.Setenv("GV_MODE", "dev")
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		e:   e,
		cfg: cfg,
	}
}
