package gv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
)

var (
	downloadedFiles = make(map[string]bool)
	mutex           = &sync.Mutex{}
	client          = resty.New().SetBaseURL("https://esm.sh")
)

// Download a file from the URL and save it to the specified destination
func (m *Manager) downloadFile(url string, dst string) error {
	// Create the output directory if it doesn't exist
	dst = strings.Replace(dst, "@types/", "", 1)
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Check if the file has already been downloaded
	mutex.Lock()
	if downloadedFiles[dst] {
		mutex.Unlock()
		_, err := os.ReadFile(dst)
		if err != nil {
			return err
		}
		return nil
	}
	downloadedFiles[dst] = true
	mutex.Unlock()

	// Fetch the content
	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to fetch %s (status: %d)", url, resp.StatusCode())
	}

	// Write to file
	if err := os.WriteFile(dst, []byte(resp.String()), 0644); err != nil {
		return err
	}

	return nil
}

type Manager struct {
	opts ManagerOptions
}

type Cdn struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Package struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type DependencyConfig struct {
	Cdns     []Cdn     `yaml:"cdns"`
	Packages []Package `yaml:"packages"`
	Types    []string  `yaml:"types"`
}

type ManagerOptions struct {
	ConfigFile string
}

func NewManager(opts ManagerOptions) *Manager {
	return &Manager{opts: opts}
}
