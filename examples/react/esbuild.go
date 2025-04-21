package main

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv/plugins"
)

var EsbuildOptions = api.BuildOptions{
	EntryPoints: []string{
		"./routes/**/*.tsx",
		"./*.html",
		"./*.ts",
		"./*.tsx",
	},
	Outdir:   "./dist",
	External: []string{"*"},
	Plugins:  []api.Plugin{plugins.CdnDependencyPlugin("./config.yaml")},
	Format:   api.FormatESModule,
	JSX:      api.JSXAutomatic,

	Bundle:            true,
	Write:             true,
	LogLevel:          api.LogLevelInfo,
	MinifySyntax:      true,
	MinifyWhitespace:  true,
	MinifyIdentifiers: true,
}
