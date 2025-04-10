package main

import (
	"log"

	"github.com/labstack/echo/v4/middleware"
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

func main() {
	plugins := []gv.Plugin{
		&plugins.ReactBabelPlugin{RootDir: "./", DistDir: "./dist"},
	}

	srv := gv.NewServer(gv.ServerConfig{
		Host:        "localhost",
		Port:        3000,
		Plugins:     plugins,
		EnableWatch: true,
	})

	group := srv.Server().Group("dist")
	group.Use(middleware.Static("./dist"))

	log.Fatal(srv.Start())
}
