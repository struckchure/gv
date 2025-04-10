package gv

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
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

type PackageMeta struct {
	Name              string                         `json:"name"`
	Version           string                         `json:"version"`
	Main              string                         `json:"main"`
	Types             string                         `json:"types"`
	TypesVersions     map[string]map[string][]string `json:"typesVersions"`
	Exports           map[string]any                 `json:"exports"`
	Dependencies      map[string]string              `json:"dependencies"`
	PeerDependencies  map[string]string              `json:"peerDependencies"`
	TypeScriptVersion string                         `json:"typeScriptVersion"`
}

func (m *Manager) Install(pkg string) error {
	pkgInfo := &PackageMeta{}
	pkgInfoRequest, err := client.R().
		SetResult(pkgInfo).
		Get(lo.Must(url.JoinPath(pkg, "package.json")))
	if err != nil {
		return err
	}

	if pkgInfoRequest.IsError() {
		return fmt.Errorf("failed to fetch %s (status: %d)", pkg, pkgInfoRequest.StatusCode())
	}

	for _, exports := range pkgInfo.Exports {
		if _, ok := exports.(map[string]any); !ok {
			continue
		}

		for _, v := range exports.(map[string]any) {
			export := v.(map[string]any)

			typePath, ok := export["types"].(string)
			if !ok {
				typePath, ok = export["default"].(string)
				if !ok {
					continue
				}
			}

			res, err := url.JoinPath(pkgInfo.Name, typePath)
			if err != nil {
				return err
			}

			color.Blue("+ %s\n", res)

			err = m.downloadFile(res, filepath.Join("types", res))
			if err != nil {
				return err
			}
		}
	}

	err = m.downloadFile(
		filepath.Join(pkgInfo.Name, "package.json"),
		filepath.Join("types", pkgInfo.Name, "package.json"),
	)
	if err != nil {
		return err
	}

	return nil
}
