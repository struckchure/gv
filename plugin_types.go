package gv

import "context"

type Context struct {
	ReqContext context.Context
	// ... add more like FS, HTTPClient, etc.
}

type ResolveResult struct {
	Id          string // Final resolved path or URL
	SideEffects bool
}

type LoadResult struct {
	Code     string
	Contents []byte // optional
	MimeType string
}

type TransformResult struct {
	Code string
	Map  string // source map (if needed)
}

type Plugin interface {
	Name() string

	OnStart() error

	// Called when resolving a module ID (bare imports, aliases, etc.)
	ResolveId(ctx *Context, id string, importer string) (*ResolveResult, error)

	// Called when loading a module after resolving
	Load(ctx *Context, id string) (*LoadResult, error)

	// Called to transform file content (e.g., JSX to JS)
	Transform(ctx *Context, code string, id string) (*TransformResult, error)

	// Called during dev HMR updates
	HandleHotUpdate(filePath string) error

	SendNotification(file string) bool
}

type PluginBase struct{}

func (p *PluginBase) Name() string {
	return "anonymous"
}

func (p *PluginBase) OnStart() error {
	return nil
}

func (p *PluginBase) ResolveId(ctx *Context, id, importer string) (*ResolveResult, error) {
	return nil, nil
}

func (p *PluginBase) Load(ctx *Context, id string) (*LoadResult, error) {
	return nil, nil
}

func (p *PluginBase) Transform(ctx *Context, code, id string) (*TransformResult, error) {
	return &TransformResult{Code: code}, nil
}

func (p *PluginBase) HandleHotUpdate(file string) error {
	return nil
}

func (p *PluginBase) SendNotification(file string) bool {
	return true
}
