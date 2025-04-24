package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
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

func CdnDependencyPlugin(depsYaml string) func(*gv.ContainerPlugin) {
	onLoad := func(args api.OnLoadArgs) (api.OnLoadResult, error) {
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

		var content string

		if args.PluginData != nil {
			previousContent, ok := args.PluginData.(api.OnLoadResult)
			if ok {
				content = *previousContent.Contents
			}
		}

		if content == "" {
			contentByte, err := os.ReadFile(args.Path)
			if err != nil {
				return api.OnLoadResult{}, err
			}

			content = string(contentByte)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return api.OnLoadResult{}, err
		}

		doc.Find("body").PrependHtml(script)

		content, err = doc.Html()
		if err != nil {
			return api.OnLoadResult{}, err
		}

		return api.OnLoadResult{
			Contents: &content,
			Loader:   api.LoaderCopy,
		}, nil
	}

	return func(cp *gv.ContainerPlugin) {
		cp.OnLoad("cdn-depenency-plugin", `\.html$`, onLoad)
	}
}
