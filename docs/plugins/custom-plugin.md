# Custom Plugin

GV supports a plugin system to customize how files are resolved, loaded, transformed, and handled during development. To write your own plugin, implement the `Plugin` interface provided by GV.

## Plugin Interface

```go
type Plugin interface {
	Name() string
	OnStart() error
	ResolveId(ctx *Context, id string, importer string) (*ResolveResult, error)
	Load(ctx *Context, id string) (*LoadResult, error)
	Transform(ctx *Context, code string, id string) (*TransformResult, error)
	HandleHotUpdate(filePath string) error
}
```

You can embed the `PluginBase` struct to inherit default no-op implementations:

```go
type MyPlugin struct {
	PluginBase
}
```

## Example Plugin

Here’s a basic example of a plugin that logs every file it transforms:

```go
type LoggerPlugin struct {
	PluginBase
}

func (p *LoggerPlugin) Name() string {
	return "logger"
}

func (p *LoggerPlugin) Transform(ctx *Context, code string, id string) (*TransformResult, error) {
	fmt.Println("Transforming:", id)
	return &TransformResult{Code: code}, nil
}
```

## Registering the Plugin

To use your plugin, add it to the list of plugins when initializing GV:

```go
plugins := []Plugin{
	&LoggerPlugin{},
	// other plugins...
}
```

## Notes

- **Name()**: should return a unique plugin name.
- **OnStart()**: is called once when the plugin system initializes.
- **ResolveId()**: is used to handle import paths (e.g., aliases, bare modules).
- **Load()**: lets you provide custom file content after resolution.
- **Transform()**: is where you can transpile or modify file contents.
- **HandleHotUpdate()**: is triggered on file change during dev (HMR).
