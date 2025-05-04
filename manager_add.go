package gv

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func (m *Manager) Add(packages ...string) error {
	cfg := DependencyConfig{}
	cfgFile, err := os.ReadFile(m.opts.ConfigFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		return err
	}

	installPackage := func(pkg string) error {
		var (
			cdn, _type, _package, _version, _path string
		)

		// Parse the package string
		parts := strings.SplitN(pkg, ":", 2)

		first := strings.SplitN(parts[0], "+", 2)
		cdn = first[0]
		if len(first) > 1 {
			_type = first[1]
		}

		last := parts[1]
		if strings.Contains(last, "/") {
			parts = strings.SplitN(last, "/", 2)
			_path = parts[1]
			last = parts[0]
		}

		parts = strings.SplitN(last, "@", 2)
		if len(parts) > 1 {
			_version = "@" + parts[1]
		}
		_package = parts[0]

		cdnUrl, ok := lo.Find(cfg.Cdns, func(_cdn Cdn) bool { return _cdn.Name == cdn })
		if !ok {
			return fmt.Errorf("CDN %s not found", cdn)
		}

		packageName := _package
		packageUrl, _ := url.JoinPath(cdnUrl.URL, fmt.Sprintf("%s%s", packageName, _version), _path)
		packageType, _ := url.JoinPath(_type, _package)

		cfg.Packages = lo.Uniq(
			append(
				cfg.Packages,
				Package{
					Name: lo.Must(url.JoinPath(packageName, _path)),
					URL:  packageUrl,
				},
			),
		)
		cfg.Types = lo.Uniq(append(cfg.Types, packageType))

		cfgFile, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		err = m.installTypes(packageType)
		if err != nil {
			return err
		}

		color.Magenta("  + %s\n", packageType)

		err = os.WriteFile(m.opts.ConfigFile, cfgFile, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	for _, _package := range packages {
		color.Green("+ %s\n", _package)
		if err := installPackage(_package); err != nil {
			color.Red("❌ Error installing %s: %v\n", _package, err)
		}
	}

	return nil
}
