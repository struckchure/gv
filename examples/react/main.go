package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
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

	log.Fatal(srv.Start())
}
