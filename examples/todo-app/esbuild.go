package main

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

var EsbuildOptions = api.BuildOptions{
	EntryPoints: []string{
		"./pages/**/*.tsx",
		"./styles/**/*.css",
		"./*.html",
		"./*.ts",
		"./*.tsx",
	},
	Outdir:   "./dist",
	External: []string{"*"},
	Plugins: []api.Plugin{
		gv.NewContainerPlugin(`\.(html|js|jsx|ts|tsx)$`).
			Setup(
				plugins.SveltePlugin,
				plugins.JsResolver,
				plugins.CdnDependencyPlugin("./config.yaml"),
			).Compose(),
	},
	Format: api.FormatESModule,
	JSX:    api.JSXAutomatic,

	Bundle:            true,
	Write:             true,
	LogLevel:          api.LogLevelInfo,
	MinifySyntax:      true,
	MinifyWhitespace:  true,
	MinifyIdentifiers: true,
}
