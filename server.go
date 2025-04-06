package gv

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	s.e.GET("/*", func(c echo.Context) error {
		requestPath := c.Request().URL.Path

		// Create GV context
		ctx := &Context{
			ReqContext: c.Request().Context(),
		}

		// Step 1: Resolve
		resolvedId, err := s.pluginContainer.ResolveId(ctx, requestPath, "")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to resolve module: "+err.Error())
		}

		// Step 2: Load
		loaded, err := s.pluginContainer.Load(ctx, resolvedId)
		if err != nil || loaded == nil {
			return c.String(http.StatusInternalServerError, "Failed to load module: "+err.Error())
		}

		// Step 3: Transform
		transformed, err := s.pluginContainer.Transform(ctx, loaded.Code, resolvedId)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Transform failed: "+err.Error())
		}

		// Step 4: Serve
		mime := loaded.MimeType
		if mime == "" {
			mime = mimeFromExt(path.Ext(resolvedId))
		}

		return c.Blob(http.StatusOK, mime, []byte(transformed.Code))
	})
}

func (s *Server) Start() error {
	s.setupRoutes()

	return s.e.Start(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
}

type ServerConfig struct {
	Host    string
	Port    int
	Plugins []Plugin
}

// NewServer sets up Echo + plugin container
func NewServer(cfg ServerConfig) *Server {
	e := echo.New()
	e.Use(middleware.CORS())

	// Create plugin container from config
	container := NewPluginContainer(cfg.Plugins...)

	server := &Server{
		e:               e,
		cfg:             cfg,
		pluginContainer: container,
	}

	return server
}
