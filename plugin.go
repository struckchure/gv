package gv

type PluginContainer struct {
	plugins []Plugin
}

func NewPluginContainer(plugins ...Plugin) *PluginContainer {
	return &PluginContainer{plugins}
}

func (pc *PluginContainer) OnStart() error {
	for _, plugin := range pc.plugins {
		err := plugin.OnStart()
		if err != nil {
			return err
		}
	}
	return nil
}

func (pc *PluginContainer) ResolveId(ctx *Context, id, importer string) (string, error) {
	for _, plugin := range pc.plugins {
		if result, err := plugin.ResolveId(ctx, id, importer); result != nil || err != nil {
			return result.Id, err
		}
	}
	return id, nil // fallback
}

func (pc *PluginContainer) Load(ctx *Context, id string) (*LoadResult, error) {
	for _, plugin := range pc.plugins {
		if result, err := plugin.Load(ctx, id); result != nil || err != nil {
			return result, err
		}
	}
	return nil, nil // fallback
}

func (pc *PluginContainer) Transform(ctx *Context, code string, id string) (*TransformResult, error) {
	result := &TransformResult{Code: code}

	for _, plugin := range pc.plugins {
		if pluginResult, err := plugin.Transform(ctx, result.Code, id); pluginResult != nil || err != nil {
			if err != nil {
				return nil, err
			}
			result = pluginResult
		}
	}
	return result, nil
}

func (pc *PluginContainer) HandleHotUpdate(file string) error {
	for _, plugin := range pc.plugins {
		if err := plugin.HandleHotUpdate(file); err != nil {
			return err
		}
	}
	return nil
}

func (pc *PluginContainer) SendNotification(file string) bool {
	for _, plugin := range pc.plugins {
		ok := plugin.SendNotification(file)
		if !ok {
			return false
		}
	}

	return true
}
