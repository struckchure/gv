package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/struckchure/gv"
)

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host:           "localhost",
		Port:           3000,
		EsBuildOptions: EsbuildOptions,
	})

	srv.Server().Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5: true,
		Root:  "dist",
	}))

	srv.Start()
}
