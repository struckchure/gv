package main

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
)

func PluginOne(args api.OnLoadArgs) (api.OnLoadResult, error) {
	var contents string

	res, ok := args.PluginData.(api.OnLoadResult)
	if ok {
		contents = *res.Contents
	} else {
		output, err := os.ReadFile(args.Path)
		if err != nil {
			return api.OnLoadResult{}, err
		}
		contents = string(output)
	}

	contents += "\nPLUGIN-ONE"

	return api.OnLoadResult{Contents: &contents, Loader: api.LoaderCopy}, nil
}

func PluginTwo(args api.OnLoadArgs) (api.OnLoadResult, error) {
	var contents string

	res, ok := args.PluginData.(api.OnLoadResult)
	if ok {
		contents = *res.Contents
	} else {
		output, err := os.ReadFile(args.Path)
		if err != nil {
			return api.OnLoadResult{}, err
		}
		contents = string(output)
	}

	contents += "\nPLUGIN-TWO"

	return api.OnLoadResult{Contents: &contents, Loader: api.LoaderCopy}, nil
}

var composedPlugins = gv.NewContainerPlugin(`\.*$`).
	OnLoad("plugin-one", `\.html$`, PluginOne).
	OnLoad("plugin-two", `\.html$`, PluginTwo).
	Compose()

func main() {
	api.Build(api.BuildOptions{
		EntryPoints: []string{"./index.html"},
		Outfile:     "./index.go.html",
		Plugins:     []api.Plugin{composedPlugins},
		LogLevel:    api.LogLevelDebug,
		Bundle:      true,
		Write:       true,
	})
}
