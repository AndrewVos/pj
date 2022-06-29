package cmd

import (
	"github.com/AndrewVos/pj/modules"
	"github.com/AndrewVos/pj/utils"
	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/mapstructure"
	"github.com/schollz/progressbar/v3"
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
		Applyables []modules.Applyable
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

	bar := progressbar.NewOptions(
		applyableCount,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, m := range modules {
		bar.Describe(m.Name)
		for _, applyable := range m.Applyables {
			err = applyable.Apply()
			if err != nil {
				return err
			}

			bar.Add(1)
		}
	}

	return err
}

func decodeApplyable(modulePath string, topLevelModule map[string]map[string]interface{}) (modules.Applyable, error) {
	if data, ok := topLevelModule["directory"]; ok {
		var applyable modules.Directory
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["symlink"]; ok {
		applyable := modules.Symlink{ModulePath: modulePath}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["group"]; ok {
		var applyable modules.Group
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["script"]; ok {
		var applyable modules.Script
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["service"]; ok {
		var applyable modules.Service
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["pacman"]; ok {
		var applyable modules.Pacman
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["aur"]; ok {
		var applyable modules.Aur
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	} else if data, ok := topLevelModule["brew"]; ok {
		var applyable modules.Brew
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
		err := mapstructure.Decode(data, &applyable)
		return applyable, err
	}

	return nil, nil
}
