package gv

import (
	_ "embed"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/fatih/color"
	"github.com/samber/lo"
)

var (
	//go:embed hmr_prefix.js
	prefixJs string

	//go:embed hmr_runtime.js
	runtimeJs string
)

func Hmr(rootHtml string) func(*ContainerPlugin) {
	loadHtml := func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var content string

		if args.PluginData != nil {
			previousContent, ok := args.PluginData.(api.OnLoadResult)
			if ok {
				content = *previousContent.Contents
			}
		}

		if content == "" {
			htmlBytes, err := os.ReadFile(rootHtml)
			if err != nil {
				return api.OnLoadResult{}, err
			}
			content = string(htmlBytes)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return api.OnLoadResult{}, err
		}

		script := fmt.Sprintf("<script>%s</script>", runtimeJs)
		doc.Find("body").PrependHtml(script)

		finalHTML, err := doc.Html()
		if err != nil {
			return api.OnLoadResult{}, err
		}

		return api.OnLoadResult{Contents: &finalHTML, Loader: api.LoaderCopy}, nil
	}

	return func(cp *ContainerPlugin) {
		cp.OnLoad("hmr-plugin", `\.(html|jsx|tsx)$`, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
			if filepath.Ext(args.Path) == ".html" {
				return loadHtml(args)
			}

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

			id, err := filepath.Rel(lo.Must(os.Getwd()), args.Path)
			if err != nil {
				return api.OnLoadResult{}, err
			}

			content = string(content) + "\n" + strings.ReplaceAll(prefixJs, "$id$", fmt.Sprintf(`"%s"`, id))
			isJsx := strings.HasSuffix(args.Path, ".jsx")

			return api.OnLoadResult{Contents: &content, Loader: lo.Ternary(isJsx, api.LoaderJSX, api.LoaderTSX)}, nil
		})

		cp.OnHMR(func(id string) (*HmrResult, error) {
			return &HmrResult{Type: HmrReloadType, Path: id}, nil
		})
	}
}

type HmrType string

const (
	HmrReloadType HmrType = "reload"
	HmrUpdateType HmrType = "update"
)

type HmrResult struct {
	Type HmrType `json:"type"`
	Path string  `json:"path"`
}

type HmrCallback func(string) (*HmrResult, error)

type HmrOptions struct {
	RootHtml          string
	WatchPath         string
	WatchExcludePaths []string
}

var HmrBroadcast = make(chan string)

func (pc *ContainerPlugin) listen() {
	for payload := range HmrBroadcast {
		for _, cb := range pc.onHmrCallbacks {
			go func() {
				res, err := cb(payload)
				if err != nil {
					color.Yellow(err.Error())
					return
				}

				if res == nil {
					return
				}

				cbPayload, err := json.Marshal(res)
				if err != nil {
					color.Yellow(err.Error())
					return
				}

				HmrClientBroadcast <- string(cbPayload)
			}()
		}
	}
}

func (pc *ContainerPlugin) WithHmr(opts HmrOptions) *ContainerPlugin {
	if os.Getenv("GV_MODE") == "dev" {
		pc.Setup(Hmr("./index.html"))

		watcher, err := NewWatcher(
			opts.WatchPath,
			opts.WatchExcludePaths,
			func(path string) { HmrBroadcast <- path },
		)
		if err != nil {
			log.Panicln("❌ Failed to create watcher:", err)
		}
		watcher.Start()

		go pc.listen()
	}

	return pc
}
