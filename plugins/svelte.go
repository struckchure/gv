package plugins

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/struckchure/gv"
)

//go:embed svelte-compiler.js
var svelteCompiler string

func SveltePlugin(cp *gv.ContainerPlugin) {
	vm := goja.New()

	vm.Set("performance", map[string]any{
		"now": func() float64 {
			return float64(time.Now().UnixNano()) / 1e6
		},
	})

	_, err := vm.RunString(svelteCompiler)
	if err != nil {
		panic(err)
	}

	svelte := vm.Get("svelte").ToObject(vm)
	compileFn := svelte.Get("compile")

	cp.OnResolve("svelte-plugin", `\.svelte$`, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		// Normalize path and change extension to .js
		ext := filepath.Ext(args.Path)
		if ext == ".svelte" {
			// Change file extension to .js and lowercase the file name
			normalizedPath := strings.TrimSuffix(args.Path, ext) + ".js"
			return api.OnResolveResult{
				Path:     normalizedPath,
				External: true,
			}, nil
		}
		return api.OnResolveResult{}, nil
	})

	cp.OnLoad("svelte-plugin", `\.svelte$`, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		source, err := os.ReadFile(args.Path)
		if err != nil {
			return api.OnLoadResult{}, err
		}

		compile, ok := goja.AssertFunction(compileFn)
		if !ok {
			return api.OnLoadResult{}, fmt.Errorf("svelte.compile is not a function")
		}

		jsResult, err := compile(goja.Undefined(), vm.ToValue(string(source)))
		if err != nil {
			return api.OnLoadResult{}, err
		}

		resultObj := jsResult.ToObject(vm)
		js := resultObj.Get("js").ToObject(vm)
		code := js.Get("code").String()

		return api.OnLoadResult{
			Contents:   &code,
			Loader:     api.LoaderJS,
			ResolveDir: filepath.Dir(args.Path),
		}, nil
	})
}
