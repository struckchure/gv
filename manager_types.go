package gv

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/samber/lo"
)

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

func (m *Manager) installTypes(types ...string) error {
	installType := func(_type string) error {
		pkgInfo := &PackageMeta{}
		pkgInfoRequest, err := client.R().
			SetResult(pkgInfo).
			Get(lo.Must(url.JoinPath(_type, "package.json")))
		if err != nil {
			return err
		}

		if pkgInfoRequest.IsError() {
			return fmt.Errorf("failed to fetch %s (status: %d)", _type, pkgInfoRequest.StatusCode())
		}

		for _, exports := range pkgInfo.Exports {
			if _, ok := exports.(map[string]any); !ok {
				continue
			}

			for _, v := range exports.(map[string]any) {
				export, ok := v.(map[string]any)
				if !ok {
					continue
				}

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

	for _, _type := range types {
		if err := installType(_type); err != nil {
			color.Red("❌ Error installing %s: %v\n", _type, err)
		}
	}

	return nil
}
