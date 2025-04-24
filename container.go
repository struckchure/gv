package gv

import (
	"path/filepath"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"
)

type ResolveCallback func(api.OnResolveArgs) (api.OnResolveResult, error)
type LoadCallback func(api.OnLoadArgs) (api.OnLoadResult, error)

type ContainerPlugin struct {
	filter string

	onResolveCallbacks map[string]func() (string, ResolveCallback)
	onLoadCallbacks    map[string]func() (string, LoadCallback)
	onHmrCallbacks     []HmrCallback
}

func (pc *ContainerPlugin) Setup(modifier ...func(*ContainerPlugin)) *ContainerPlugin {
	for _, modifier := range modifier {
		modifier(pc)
	}

	return pc
}

func (pc *ContainerPlugin) OnResolve(name, filter string, callback ResolveCallback) *ContainerPlugin {
	pc.onResolveCallbacks[name] = func() (string, ResolveCallback) { return filter, callback }

	return pc
}

func (pc *ContainerPlugin) OnLoad(name, filter string, callback LoadCallback) *ContainerPlugin {
	pc.onLoadCallbacks[name] = func() (string, LoadCallback) { return filter, callback }

	return pc
}

func (pc *ContainerPlugin) OnHMR(callback HmrCallback) *ContainerPlugin {
	pc.onHmrCallbacks = append(pc.onHmrCallbacks, callback)

	return pc
}

func (pc *ContainerPlugin) Compose() api.Plugin {
	return api.Plugin{
		Name: "plugin-container",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: pc.filter}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				result := api.OnResolveResult{}

				for _, plugins := range pc.onResolveCallbacks {
					filter, callback := plugins()
					pattern, err := regexp.MatchString(filter, filepath.Clean(args.Path))
					if err != nil {
						return api.OnResolveResult{}, err
					}
					if !pattern {
						continue
					}

					_result, err := callback(args)
					if err != nil {
						return api.OnResolveResult{}, err
					}

					result = _result

					args.PluginData = _result
				}

				return result, nil
			})

			build.OnLoad(api.OnLoadOptions{Filter: pc.filter}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				result := api.OnLoadResult{}

				for _, plugins := range pc.onLoadCallbacks {
					filter, callback := plugins()
					pattern, err := regexp.MatchString(filter, filepath.Clean(args.Path))
					if err != nil {
						return api.OnLoadResult{}, err
					}
					if !pattern {
						continue
					}

					_result, err := callback(args)
					if err != nil {
						return api.OnLoadResult{}, err
					}

					result = _result

					args.PluginData = _result
				}

				return result, nil
			})
		},
	}
}

func NewContainerPlugin(filter string) *ContainerPlugin {
	return &ContainerPlugin{
		filter: filter,

		onResolveCallbacks: map[string]func() (string, ResolveCallback){},
		onLoadCallbacks:    map[string]func() (string, LoadCallback){},
		onHmrCallbacks:     []HmrCallback{},
	}
}
