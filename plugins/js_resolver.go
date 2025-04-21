package plugins

import (
	"github.com/evanw/esbuild/pkg/api"
)

func JsResolver() api.Plugin {
	return api.Plugin{
		Name:  "js-resolver",
		Setup: func(build api.PluginBuild) {},
	}
}
