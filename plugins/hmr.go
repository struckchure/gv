package plugins

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/struckchure/gv"
)

//go:embed hmr.js
var hmrJS string

type HMRPlugin struct {
	gv.PluginBase
}

func (f *HMRPlugin) Name() string {
	return "gv/hmr-plugin"
}

func (f *HMRPlugin) Transform(ctx *gv.Context, code, id string) (*gv.TransformResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(code))
	if err != nil {
		return nil, err
	}

	script := fmt.Sprintf("<script>%s</script>", hmrJS)
	doc.Find("body").AppendHtml(script)

	code, err = doc.Html()
	if err != nil {
		return nil, err
	}

	return &gv.TransformResult{Code: code}, nil
}
