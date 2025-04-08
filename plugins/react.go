package plugins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	babel "github.com/jvatic/goja-babel"
	"github.com/struckchure/gv"
)

type ReactPlugin struct {
	gv.PluginBase
	RootDir string
	DistDir string
}

func (f *ReactPlugin) Name() string {
	return "react-plugin"
}

// Initialize with default values if needed
func (f *ReactPlugin) build() error {
	// === Recursively walk through all files ===
	err := filepath.WalkDir(f.RootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the output directory
		if strings.HasPrefix(path, f.DistDir) {
			return nil
		}

		// ignore dist directory
		if d.IsDir() && d.Name() == "dist" {
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".jsx") {
			return f.buildJSX(path, filepath.Dir(path), filepath.Join(f.DistDir, filepath.Dir(path)))
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".js") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = os.WriteFile(filepath.Join(f.DistDir, path), content, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *ReactPlugin) buildJSX(path, inputDir, outputDir string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", path, err)
	}

	res, err := babel.Transform(strings.NewReader(string(content)), map[string]interface{}{
		"plugins": []string{
			"transform-react-jsx",
			"transform-block-scoping",
		},
	})
	if err != nil {
		return fmt.Errorf("babel transform failed for %s: %w", path, err)
	}

	// Compute the output path
	relPath, err := filepath.Rel(inputDir, path)
	if err != nil {
		return err
	}
	jsOutputPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".jsx")+".js")

	// Create the parent directory if it doesn’t exist
	if err := os.MkdirAll(filepath.Dir(jsOutputPath), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(jsOutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, res); err != nil {
		return err
	}

	fmt.Println(color.BlueString(path) + color.GreenString(" → %s", jsOutputPath))

	return nil
}

func (f *ReactPlugin) OnStart() error {
	return f.build()
}

func (f *ReactPlugin) ResolveId(ctx *gv.Context, id, importer string) (*gv.ResolveResult, error) {
	return &gv.ResolveResult{Id: filepath.Clean(id)}, nil
}

func (f *ReactPlugin) Load(ctx *gv.Context, fullPath string) (*gv.LoadResult, error) {
	// Initialize plugin if needed
	err := f.build()
	if err != nil {
		return nil, err
	}

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

func (f *ReactPlugin) HandleHotUpdate(file string) error {
	if filepath.Ext(file) != ".jsx" {
		return nil
	}

	return f.buildJSX(file, filepath.Dir(file), filepath.Join(f.DistDir, filepath.Dir(file)))
}
