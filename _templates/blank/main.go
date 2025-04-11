package main

import "github.com/struckchure/gv"

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host:    "0.0.0.0",
		Port:    3000,
		Plugins: []gv.Plugin{},
	})
	srv.Start()
}
