package plugins

import (
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
)

// TODO: resolve default to `index.js` if exists
// e.g import Mod from "./page"; // where `page` dir has index.js
func JsResolver(cp *gv.ContainerPlugin) {
	cp.OnResolve("js-resolver", `.+`, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		ext := filepath.Ext(args.Path)
		if strings.HasPrefix(args.Path, ".") && ext != ".js" {
			args.Path += ".js"
		}

		return api.OnResolveResult{
			Path:     args.Path,
			External: true,
		}, nil
	})
}
