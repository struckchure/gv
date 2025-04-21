package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/evanw/esbuild/pkg/api"
	"gopkg.in/yaml.v3"
)

type CdnDependencyPackage struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type CdnDependencyConfig struct {
	Packages []CdnDependencyPackage `yaml:"packages"`
}

type importmap struct {
	Imports map[string]string `json:"imports"`
}

func CdnDependencyPlugin(depsYaml string) api.Plugin {
	return api.Plugin{
		Name: "cdn-dependency-plugin",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: `\.html$`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				depsContent, err := os.ReadFile(depsYaml)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				var depsConfig CdnDependencyConfig
				if err := yaml.Unmarshal(depsContent, &depsConfig); err != nil {
					return api.OnLoadResult{}, err
				}

				im := importmap{Imports: make(map[string]string)}
				for _, pkg := range depsConfig.Packages {
					im.Imports[pkg.Name] = pkg.URL
				}

				mapContent, err := json.Marshal(im)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				script := fmt.Sprintf(`<script type="importmap">%s</script>`, string(mapContent))

				htmlBytes, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBytes)))
				if err != nil {
					return api.OnLoadResult{}, err
				}

				doc.Find("body").PrependHtml(script)

				finalHTML, err := doc.Html()
				if err != nil {
					return api.OnLoadResult{}, err
				}

				return api.OnLoadResult{
					Contents:   &finalHTML,
					Loader:     api.LoaderCopy,
					ResolveDir: filepath.Dir(args.Path),
				}, nil
			})
		},
	}
}
