package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
)

type ReactEsBuildPlugin struct {
	gv.PluginBase
	RootDir     string
	DistDir     string
	EntryPoints []string
}

var extensions = []string{".jsx", ".js", ".tsx", ".ts", ".css"}

func (f *ReactEsBuildPlugin) transpileDirectory() error {
	result := api.Build(api.BuildOptions{
		EntryPoints: f.EntryPoints,
		Outdir:      "dist",
		External:    []string{"*"},
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
		Format:      api.FormatESModule,
		Banner: map[string]string{
			"js": `import React from "react";`,
		},
	})

	if len(result.Errors) > 0 {
		return errors.New(result.Errors[0].Text)
	}

	return nil
}

func (f *ReactEsBuildPlugin) transpileFile(path, outputDir string) error {
	ext := filepath.Ext(path)
	if lo.Contains([]string{".tsx", ".jsx", ".ts", ".js"}, ext) {
		path = strings.Replace(filepath.Base(path), filepath.Ext(path), ".js", 1)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Outfile:     filepath.Join(outputDir, path),
		External:    []string{"*"},
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
		Format:      api.FormatESModule,
		Banner: map[string]string{
			"js": `import React from "react";`,
		},
	})

	if len(result.Errors) > 0 {
		return errors.New(result.Errors[0].Text)
	}

	return nil
}

func (f *ReactEsBuildPlugin) Name() string {
	return "gv/react-esbuild-plugin"
}

func (f *ReactEsBuildPlugin) OnStart() error {
	return f.transpileDirectory()
}

func (f *ReactEsBuildPlugin) ResolveId(ctx *gv.Context, id, importer string) (*gv.ResolveResult, error) {
	return &gv.ResolveResult{Id: filepath.Clean(id)}, nil
}

func (f *ReactEsBuildPlugin) Load(ctx *gv.Context, fullPath string) (*gv.LoadResult, error) {
	rootHtml, err := os.ReadFile("./index.html")
	if err != nil {
		return nil, err
	}

	return &gv.LoadResult{
		Contents: rootHtml,
		Code:     string(rootHtml),
		MimeType: "text/html",
	}, nil
}

func (f *ReactEsBuildPlugin) HandleHotUpdate(file string) error {
	if strings.HasPrefix(filepath.Clean(file), filepath.Clean(f.DistDir)) {
		return nil
	}

	if !lo.Contains(extensions, filepath.Ext(file)) {
		return nil
	}

	return f.transpileFile(file, filepath.Join(f.DistDir, filepath.Dir(file)))
}

func (f *ReactEsBuildPlugin) SendNotification(file string) bool {
	return strings.HasPrefix(filepath.Clean(file), filepath.Clean(f.DistDir))
}
