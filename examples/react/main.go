package main

import (
	"embed"
	"log"

	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
)

//go:embed dist
var _ embed.FS

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host: "localhost",
		Port: 3000,

		EsBuildOptions:    EsbuildOptions,
		WatchPath:         lo.ToPtr("./"),
		WatchExcludePaths: &[]string{"dist"},
	})

	srv.Server().Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5: true,
		Root:  "dist",
	}))

	log.Fatal(srv.Start())
}
