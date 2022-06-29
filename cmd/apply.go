package cmd

import (
	"fmt"
	"github.com/AndrewVos/pj/applyables"
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
		Name       string
		Applyables []applyables.Applyable
	}

	modules := []Module{}
	applyableCount := 0

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
			applyable, err := decodeApplyable(modulePath, topLevelModule)
			if err != nil {
				return err
			}
			m.Applyables = append(m.Applyables, applyable)
			applyableCount += 1
		}

		modules = append(modules, m)
	}

	for _, m := range modules {
		if Verbose {
			fmt.Printf("Applying module %s...\n", m.Name)
		}
		for _, applyable := range m.Applyables {
			err = applyable.Apply()
			if err != nil {
				return err
			}
		}
	}

	return err
}

func decodeApplyable(modulePath string, topLevelModule map[string]map[string]interface{}) (applyables.Applyable, error) {
	if data, ok := topLevelModule["directory"]; ok {
		var applyable applyables.Directory
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["symlink"]; ok {
		applyable := applyables.Symlink{ModulePath: modulePath}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["group"]; ok {
		var applyable applyables.Group
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["script"]; ok {
		var applyable applyables.Script
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["service"]; ok {
		var applyable applyables.Service
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["pacman"]; ok {
		var applyable applyables.Pacman
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["aur"]; ok {
		var applyable applyables.Aur
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["brew"]; ok {
		var applyable applyables.Brew
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	}

	return nil, nil
}
