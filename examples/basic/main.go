package main

import "github.com/struckchure/gv"

func main() {
	srv := gv.NewServer(gv.ServerConfig{})
	srv.Start()
}
