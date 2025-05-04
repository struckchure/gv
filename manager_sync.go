package gv

import (
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func (m *Manager) Sync(depsFile string) {
	content, err := os.ReadFile(depsFile)
	if err != nil {
		color.Red(err.Error())
		return
	}

	config := &DependencyConfig{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		color.Red(err.Error())
		return
	}

	if err := m.installTypes(config.Types...); err != nil {
		color.Red(err.Error())
		return
	}

	color.Green("Sync Complete.")
}
