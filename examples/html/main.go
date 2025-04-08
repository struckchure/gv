package main

import (
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

func main() {
	plugins := []gv.Plugin{
		&plugins.HTMLPlugin{RootDir: ".", ChildrenSelector: "#children"},
	}

	srv := gv.NewServer(gv.ServerConfig{
		Host:    "localhost",
		Port:    3000,
		Plugins: plugins,
	})
	srv.Start()
}
