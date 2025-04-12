package gv

type Manager struct{}

type Package struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type DependencyConfig struct {
	Packages []Package `yaml:"packages"`
	Types    []string  `yaml:"types"`
}

func NewManager() *Manager {
	return &Manager{}
}
