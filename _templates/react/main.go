package main

import (
	"log"

	"github.com/labstack/echo/v4/middleware"
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

func main() {
	plugins := []gv.Plugin{
		&plugins.ReactEsBuildPlugin{
			RootDir:     "./",
			DistDir:     "./dist",
			EntryPoints: []string{"./**/*.tsx", "index.css"},
		},
		&plugins.CdnDepencyPlugin{DepsYaml: "./deps.yaml"},
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
