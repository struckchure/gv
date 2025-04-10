package gv

import (
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
	go func() {
		defer close(done)
		for {
			select {
			case msg, ok := <-eventChan:
				if !ok {
					return // Channel closed, exit goroutine
				}

				if err := s.pluginContainer.HandleHotUpdate(msg); err != nil {
					log.Println(err)
				}
			case <-done:
				return
			}
		}
	}()

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
