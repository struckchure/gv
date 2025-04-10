package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/struckchure/gv"
	"gopkg.in/yaml.v3"
)

type CdnDepencyPlugin struct {
	gv.PluginBase
	RootDir string

	DepsYaml string
}

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

func (f *CdnDepencyPlugin) Name() string {
	return "gv/cdn-depency-plugin"
}

func (f *CdnDepencyPlugin) Transform(ctx *gv.Context, code string, id string) (*gv.TransformResult, error) {
	depsConfig := &CdnDependencyConfig{}

	depsContent, err := os.ReadFile(f.DepsYaml)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(depsContent, &depsConfig)
	if err != nil {
		return nil, err
	}

	im := importmap{
		Imports: map[string]string{},
	}

	for _, pkg := range depsConfig.Packages {
		im.Imports[pkg.Name] = pkg.URL
	}

	content, err := json.Marshal(&im)
	if err != nil {
		return nil, err
	}

	script := fmt.Sprintf(`<script type="importmap">%s</script>`, string(content))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(code))
	if err != nil {
		return nil, err
	}

	doc.Find("body").PrependHtml(script)

	code, err = doc.Html()
	if err != nil {
		return nil, err
	}

	return &gv.TransformResult{Code: code}, nil
}
