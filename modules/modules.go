package modules

import (
	"github.com/AndrewVos/pj/actions"
	"github.com/AndrewVos/pj/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
)

type Module struct {
	Name    string
	Path    string
	Actions []actions.Action
}

func (m Module) Apply() error {
	for _, action := range m.Actions {
		err := action.Apply(m.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func modulePaths() ([]string, error) {
	var modulePaths []string

	cwd, err := os.Getwd()
	if err != nil {
		return modulePaths, err
	}

	paths, err := filepath.Glob(filepath.Join(cwd, "modules", "*"))
	if err != nil {
		return modulePaths, err
	}

	for _, path := range paths {
		exists, err := utils.DirectoryExists(path)
		if err != nil {
			return modulePaths, err
		}

		if exists {
			modulePaths = append(modulePaths, path)
		}
	}
	return modulePaths, nil
}

func LoadModules() ([]Module, error) {
	var modules []Module

	modulePaths, err := modulePaths()
	if err != nil {
		return modules, err
	}

	for _, modulePath := range modulePaths {
		m := Module{Name: path.Base(modulePath), Path: modulePath}

		configurationPath := filepath.Join(modulePath, "configuration.yml")

		if utils.FileExists(configurationPath) {
			contents, err := utils.ReadFile(configurationPath)
			if err != nil {
				return modules, err
			}

			document := []map[string]map[string]interface{}{}
			err = yaml.Unmarshal([]byte(contents), &document)
			if err != nil {
				return modules, err
			}

			for _, topLevelModule := range document {
				action, err := decodeAction(modulePath, topLevelModule)
				if err != nil {
					return modules, err
				}
				m.Actions = append(m.Actions, action)
			}
		}

		modules = append(modules, m)
	}

	return modules, nil
}

func findActionFor(topLevelModule map[string]map[string]interface{}) (actions.Action, map[string]interface{}) {
	for _, action := range actions.All {
		if data, ok := topLevelModule[action.Flag()]; ok {
			return action, data
		}
	}
	panic("action not found")
}

func decodeAction(modulePath string, topLevelModule map[string]map[string]interface{}) (actions.Action, error) {
	action, data := findActionFor(topLevelModule)

	if _, ok := action.(actions.Pacman); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	} else if _, ok := action.(actions.Aur); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	} else if _, ok := action.(actions.Brew); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	}

	err := mapstructure.Decode(data, &action)
	return action, err
}
