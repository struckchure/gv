package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
	index "github.com/struckchure/gv/examples/todo-app/pages"
)

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host: "localhost",
		Port: 3000,

		EsBuildOptions: EsbuildOptions,
	})

	e := srv.Server()

	e.GET("/*", func(c echo.Context) error {
		content, err := os.ReadFile(filepath.Join(lo.Must(os.Getwd()), "/dist/index.html"))
		if err != nil {
			return err
		}

		return c.HTML(200, string(content))
	})
	g := e.Group("/dist")
	g.Static("/", "dist")

	api := e.Group("api")
	api.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	todoPage := index.Page{}
	api.GET("/todos/", todoPage.List)
	api.POST("/todos/", todoPage.Create)
	api.PATCH("/todos/:id/", todoPage.Update)
	api.GET("/todos/:id/", todoPage.Get)
	api.DELETE("/todos/:id/", todoPage.Delete)

	log.Fatal(srv.Start())
}
