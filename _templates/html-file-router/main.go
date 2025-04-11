package main

import (
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

func main() {
	srv := gv.NewServer(gv.ServerConfig{
		Host:    "0.0.0.0",
		Port:    3000,
		Plugins: []gv.Plugin{&plugins.HTMLPlugin{RootDir: ".", ChildrenSelector: "#children"}},
	})
	srv.Start()
}
