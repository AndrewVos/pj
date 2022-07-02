package cmd

import (
	"fmt"
	"github.com/AndrewVos/pj/actions"
	"github.com/AndrewVos/pj/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"path/filepath"
)

var Verbose bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := apply()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	applyCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(applyCmd)
}

func apply() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	modulePaths, err := filepath.Glob(filepath.Join(cwd, "modules", "*"))
	if err != nil {
		return err
	}

	type Module struct {
		Name    string
		Actions []actions.Action
	}

	modules := []Module{}
	actionCount := 0

	for _, modulePath := range modulePaths {
		m := Module{Name: path.Base(modulePath)}

		configurationPath := filepath.Join(modulePath, "configuration.yml")

		contents, err := utils.ReadFile(configurationPath)
		if err != nil {
			return err
		}

		document := []map[string]map[string]interface{}{}
		err = yaml.Unmarshal([]byte(contents), &document)
		if err != nil {
			return err
		}

		for _, topLevelModule := range document {
			action, err := decodeAction(modulePath, topLevelModule)
			if err != nil {
				return err
			}
			m.Actions = append(m.Actions, action)
			actionCount += 1
		}

		modules = append(modules, m)
	}

	for _, m := range modules {
		if Verbose {
			fmt.Printf("Applying module %s...\n", m.Name)
		}
		for _, action := range m.Actions {
			err = action.Apply()
			if err != nil {
				return err
			}
		}
	}

	return err
}

func decodeAction(modulePath string, topLevelModule map[string]map[string]interface{}) (actions.Action, error) {
	if data, ok := topLevelModule["directory"]; ok {
		var action actions.Directory
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["symlink"]; ok {
		action := actions.NewSymlink(modulePath)
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["group"]; ok {
		var action actions.Group
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["script"]; ok {
		var action actions.Script
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["service"]; ok {
		var action actions.Service
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["pacman"]; ok {
		var action actions.Pacman
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["aur"]; ok {
		var action actions.Aur
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &action)
		return action, err
	} else if data, ok := topLevelModule["brew"]; ok {
		var action actions.Brew
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &action)
		return action, err
	}

	return nil, nil
}
