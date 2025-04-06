package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/struckchure/gv"
)

type ReactPlugin struct {
	gv.PluginBase
	RootDir string

	// default: `#children`
	RootSelector string
}

func (f *ReactPlugin) Name() string {
	return "html-router"
}

// Initialize with default values if needed
func (f *ReactPlugin) init() {
	if f.RootSelector == "" {
		f.RootSelector = "#children"
	}
}

func (f *ReactPlugin) ResolveId(ctx *gv.Context, id, importer string) (*gv.ResolveResult, error) {
	// Clean and resolve path relative to the root
	cleanPath := filepath.Clean(id)

	// First, try exact match for +page.html
	fullPath := filepath.Join(f.RootDir, cleanPath, "+page.html")
	info, err := os.Stat(fullPath)
	if err == nil && !info.IsDir() {
		return &gv.ResolveResult{Id: fullPath}, nil
	}

	// Then, try with .html extension if not already present
	if !strings.HasSuffix(cleanPath, ".html") {
		fullPathWithExt := filepath.Join(f.RootDir, cleanPath+".html")
		info, err := os.Stat(fullPathWithExt)
		if err == nil && !info.IsDir() {
			return &gv.ResolveResult{Id: fullPathWithExt}, nil
		}
	}

	// Let other plugins try
	return nil, nil
}

func (f *ReactPlugin) Load(ctx *gv.Context, id string) (*gv.LoadResult, error) {
	// Initialize plugin if needed
	f.init()

	// Verify path is within the root directory
	fullPath := id
	if !strings.HasPrefix(fullPath, f.RootDir) {
		fullPath = filepath.Join(f.RootDir, id)
	}

	// Check if file exists
	stat, err := os.Stat(fullPath)
	if err != nil || stat.IsDir() {
		return nil, errors.New("file not found")
	}

	return &gv.LoadResult{}, nil
}
