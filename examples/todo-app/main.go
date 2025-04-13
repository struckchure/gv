package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/struckchure/gv"
	index "github.com/struckchure/gv/examples/todo-app/pages"
	"github.com/struckchure/gv/plugins"
)

func main() {
	plugins := []gv.Plugin{
		&plugins.ReactEsBuildPlugin{
			RootDir:     "./",
			DistDir:     "./dist",
			EntryPoints: []string{"./**/*.tsx", "./*.ts", "./styles/**/*.css"},
		},
		&plugins.CdnDepencyPlugin{DepsYaml: "./deps.yaml"},
		&plugins.HMRPlugin{},
	}

	srv := gv.NewServer(gv.ServerConfig{
		Host:        "localhost",
		Port:        3000,
		Plugins:     plugins,
		EnableWatch: true,
	})

	group := srv.Server().Group("dist")
	group.Use(middleware.Static("./dist"))

	api := srv.Server().Group("api")
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
