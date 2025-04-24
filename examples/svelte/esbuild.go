package main

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
	"github.com/struckchure/gv/plugins"
)

var EsbuildOptions = api.BuildOptions{
	EntryPoints: []string{
		"./*.html",
		"./*.ts",
		"./*.svelte",
	},
	Outdir:   "./dist",
	External: []string{"*"},
	Plugins: []api.Plugin{
		gv.NewContainerPlugin(`\.(html|svelte)$`).
			Setup(
				plugins.SveltePlugin,
				plugins.CdnDependencyPlugin("./config.yaml"),
			).Compose(),
	},
	Format: api.FormatESModule,

	Bundle:            true,
	Write:             true,
	LogLevel:          api.LogLevelInfo,
	MinifySyntax:      true,
	MinifyWhitespace:  true,
	MinifyIdentifiers: true,
}
