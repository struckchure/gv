package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
)

var PORT, _ = strconv.Atoi(lo.Ternary(os.Getenv("PORT") == "", "3000", os.Getenv("PORT")))

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host:           "0.0.0.0",
		Port:           PORT,
		EsBuildOptions: EsbuildOptions,
	})

	if os.Getenv("GV_MODE") == "build" {
		if err := srv.Build(); err != nil {
			log.Fatal(err)
		}
		return
	}

	build := exec.Command("ls", "-la", "dist")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	srv.Server().Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5: true,
		Root:  "dist",
	}))

	log.Fatal(srv.Start())
}
