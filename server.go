package gv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Server struct {
	e *echo.Echo

	cfg             ServerConfig
	pluginContainer *PluginContainer
}

// Helper to guess content type
func mimeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".js":
		return "application/javascript"
	case ".ts":
		return "application/javascript"
	case ".tsx":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	default:
		return "application/octet-stream"
	}
}

func (s *Server) setupRoutes() {
	s.pluginContainer.OnStart()

	eventChan := s.cfg.eventBus.Subscribe(string(FileUpdated))
	done := make(chan struct{})

	type hmrPayload struct {
		Type string `json:"type"`
		File string `json:"file"`
	}

	s.e.GET("/_/ws/", func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			defer close(done)

			// Ping/Pong mechanism
			// go func() {
			// 	ticker := time.NewTicker(10 * time.Second) // Ping interval
			// 	defer ticker.Stop()

			// 	for {
			// 		select {
			// 		case <-ticker.C:
			// 			// Send a Ping frame to the client
			// 			if err := websocket.JSON.Send(ws, map[string]string{"type": "ping"}); err != nil {
			// 				c.Logger().Error("Failed to send ping:", err)
			// 				return // Exit the goroutine on error (likely client disconnected)
			// 			}
			// 			c.Logger().Info("Sent ping to client")

			// 		case <-done:
			// 			return
			// 		}
			// 	}
			// }()

			for {
				select {
				case file, ok := <-eventChan:
					if !ok {
						return // Channel closed, exit goroutine
					}

					if err := s.pluginContainer.HandleHotUpdate(file); err != nil {
						return
					}

					shouldNotify := s.pluginContainer.SendNotification(file)
					if shouldNotify {
						payload, err := json.Marshal(&hmrPayload{
							Type: fmt.Sprintf("%s:update", strings.Replace(filepath.Ext(file), ".", "", 1)),
							File: filepath.Join("/", file),
						})
						if err != nil {
							c.Logger().Error(err)
							return
						}
						err = websocket.Message.Send(ws, string(payload))
						if err != nil {
							c.Logger().Error(err)
						}
					}
				case <-done:
					return
				}
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	})

	s.e.GET("/*", func(c echo.Context) error {
		requestPath := c.Request().URL.Path

		// Create GV context
		ctx := &Context{
			ReqContext: c.Request().Context(),
		}

		resolvedId, err := s.pluginContainer.ResolveId(ctx, requestPath, "")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to resolve module: "+err.Error())
		}

		loaded, err := s.pluginContainer.Load(ctx, resolvedId)
		if err != nil || loaded == nil {
			return c.String(http.StatusInternalServerError, "Failed to load module: "+err.Error())
		}

		transformed, err := s.pluginContainer.Transform(ctx, loaded.Code, resolvedId)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Transform failed: "+err.Error())
		}

		mime := loaded.MimeType
		if mime == "" {
			mime = mimeFromExt(path.Ext(resolvedId))
		}

		return c.Blob(http.StatusOK, mime, []byte(transformed.Code))
	})
}

func (s *Server) Server() *echo.Echo {
	return s.e
}

func (s *Server) Start() error {
	s.setupRoutes()

	fmt.Println("\n" + color.GreenString("➜") + " Local: " + color.MagentaString(fmt.Sprintf("http://%s:%d", s.cfg.Host, s.cfg.Port)) + "\n")

	if s.cfg.EnableWatch {
		go s.Watch()
	}

	return s.e.Start(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
}

type EventType string

const (
	FileCreated EventType = "file:created"
	FileUpdated EventType = "file:updated"
	FileDeleted EventType = "file:deleted"
)

func (s *Server) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					info, err := os.Stat(event.Name)
					if err != nil {
						log.Fatal(err)
						return
					}

					if info.IsDir() {
						watcher.Add(event.Name) // Watch new folder
						return
					}

					s.cfg.eventBus.Publish(string(FileCreated), event.Name)
				}

				if event.Has(fsnotify.Write) {
					s.cfg.eventBus.Publish(string(FileUpdated), event.Name)
				}

				if event.Has(fsnotify.Remove) {
					s.cfg.eventBus.Publish(string(FileDeleted), event.Name)
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	color.Green("Watching for changes...")

	// Block main goroutine forever.
	<-make(chan struct{})
}

type ServerConfig struct {
	Host        string
	Port        int
	Plugins     []Plugin
	EnableWatch bool

	eventBus *EventBus
}

// NewServer sets up Echo + plugin container
func NewServer(cfg ServerConfig) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	cfg.eventBus = NewEventBus()

	return &Server{
		e:               e,
		cfg:             cfg,
		pluginContainer: NewPluginContainer(cfg.Plugins...),
	}
}
