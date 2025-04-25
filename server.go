package gv

import (
	"errors"
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
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

var HmrClientBroadcast = make(chan string)

func (s *Server) HandleHMR() error {
	color.Magenta("[HMR] Wating for client to connect ...")

	handler := func(c echo.Context) error {
		// Upgrade connection to WebSocket
		websocket.Handler(func(ws *websocket.Conn) {
			if ws.IsClientConn() {
				color.Magenta("[HMR] Client connected")
			}

			defer ws.Close()

			for payload := range HmrClientBroadcast {
				if err := websocket.Message.Send(ws, payload); err != nil {
					c.Logger().Error("WebSocket broadcast error:", err)
					continue
				}
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	}

	s.e.GET("/__hmr__", handler)

	return nil
}

func (s *Server) Start() error {
	if os.Getenv("GV_MODE") == "dev" {
		s.Watch()
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
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		e:   e,
		cfg: cfg,
	}
}
