package plugins

import (
	"errors"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/struckchure/gv"
)

type HTMLPlugin struct {
	gv.PluginBase
	RootDir string

	// default: `#children`
	ChildrenSelector string
}

func (f *HTMLPlugin) Name() string {
	return "html-router"
}

// Initialize with default values if needed
func (f *HTMLPlugin) init() {
	if f.ChildrenSelector == "" {
		f.ChildrenSelector = "#children"
	}
}

func (f *HTMLPlugin) ResolveId(ctx *gv.Context, id, importer string) (*gv.ResolveResult, error) {
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

func (f *HTMLPlugin) applyLayout(layoutPath string, contents *string, selector string) error {
	layoutData, err := os.ReadFile(layoutPath)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(layoutData)))
	if err != nil {
		return err
	}

	contentsDom, err := goquery.NewDocumentFromReader(strings.NewReader(*contents))
	if err != nil {
		return err
	}

	doc.Find(selector).ReplaceWithHtml(lo.Must(contentsDom.Html()))
	updated, err := doc.Html()
	if err != nil {
		return err
	}

	*contents = updated
	return nil
}

// applyDirectContent reads content from an override file and replaces the current content
func (f *HTMLPlugin) applyDirectContent(overridePath string, contents *string) error {
	overrideContent, err := os.ReadFile(overridePath)
	if err != nil {
		return errors.New("failed to read override page: " + err.Error())
	}
	*contents = string(overrideContent)
	return nil
}

// getLayoutPaths returns a slice of paths from the current directory to the root
func (f *HTMLPlugin) getLayoutPaths(currentPath string) ([]string, error) {
	if currentPath == "" {
		return nil, errors.New("current path is empty")
	}

	// Make path relative to the root directory
	relPath, err := filepath.Rel(f.RootDir, currentPath)
	if err != nil {
		return nil, err
	}

	// If we're at the root, return an empty slice
	if relPath == "." {
		return []string{f.RootDir}, nil
	}

	// Split the path into segments
	segments := strings.Split(relPath, string(os.PathSeparator))

	// Build paths from innermost to outermost
	paths := make([]string, len(segments)+1)

	// Start with the innermost directory (current directory)
	paths[0] = currentPath

	// Build parent directories moving outward
	for i := 1; i < len(segments); i++ {
		paths[i] = filepath.Dir(paths[i-1])
	}

	// Add the root directory as the last entry
	paths[len(segments)] = f.RootDir

	return paths, nil
}

func (f *HTMLPlugin) Load(ctx *gv.Context, id string) (*gv.LoadResult, error) {
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

	// Read file
	byteContents, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	contents := string(byteContents)

	// Determine MIME type
	mimeType := mime.TypeByExtension(filepath.Ext(fullPath))
	if mimeType == "" {
		mimeType = "application/javascript" // default
	}

	// Check for @page.html override in the current directory
	pageDir := filepath.Dir(fullPath)
	overridePagePath := filepath.Join(pageDir, "@page.html")
	if _, err := os.Stat(overridePagePath); err == nil {
		// Use override and skip all layouts
		if err := f.applyDirectContent(overridePagePath, &contents); err != nil {
			return nil, err
		}
		return &gv.LoadResult{
			Contents: []byte(contents),
			Code:     contents,
			MimeType: mimeType,
		}, nil
	}

	// Get layout paths from current directory to root
	layoutPaths, err := f.getLayoutPaths(pageDir)
	if err != nil {
		return nil, err
	}

	// Apply layouts from innermost to outermost
	for _, dirPath := range layoutPaths {
		// Check for @layout.html (override that stops inheritance)
		overrideLayoutPath := filepath.Join(dirPath, "@layout.html")
		if _, err := os.Stat(overrideLayoutPath); err == nil {
			if err := f.applyLayout(overrideLayoutPath, &contents, f.ChildrenSelector); err != nil {
				return nil, err
			}
			// Stop layout inheritance after override
			break
		}

		// Apply regular +layout.html if it exists
		layoutPath := filepath.Join(dirPath, "+layout.html")
		if _, err := os.Stat(layoutPath); err == nil {
			if err := f.applyLayout(layoutPath, &contents, f.ChildrenSelector); err != nil {
				return nil, err
			}
		}
	}

	return &gv.LoadResult{
		Contents: []byte(contents),
		Code:     contents,
		MimeType: mimeType,
	}, nil
}
