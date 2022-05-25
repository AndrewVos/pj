package cmd

import (
	"github.com/AndrewVos/pj/modules"
	"github.com/AndrewVos/pj/utils"
	"github.com/k0kubun/go-ansi"
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
			m.Applyables = append(m.Applyables, buildModule(modulePath, topLevelModule))
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

func buildModule(modulePath string, topLevelModule map[string]map[string]interface{}) modules.Applyable {
	if module, ok := topLevelModule["pacman"]; ok {
		return retrievePacman(module)
	} else if module, ok := topLevelModule["aur"]; ok {
		return retrieveAur(module)
	} else if module, ok := topLevelModule["symlink"]; ok {
		return retrieveSymlink(modulePath, module)
	} else if module, ok := topLevelModule["group"]; ok {
		return retrieveGroup(module)
	} else if module, ok := topLevelModule["script"]; ok {
		return retrieveScript(module)
	} else if module, ok := topLevelModule["service"]; ok {
		return retrieveService(module)
	} else if module, ok := topLevelModule["directory"]; ok {
		return retrieveDirectory(module)
	}

	return nil
}

func retrievePacman(module map[string]interface{}) modules.Pacman {
	pacman := modules.Pacman{}

	if value, ok := module["name"]; ok {
		if name, ok := value.(string); ok {
			pacman.Name = []string{name}
		} else if nameValues, ok := value.([]interface{}); ok {
			names := []string{}
			for _, nameValue := range nameValues {
				if name, ok := nameValue.(string); ok {
					names = append(names, name)
				}
			}
			pacman.Name = names
		}
	}
	return pacman
}

func retrieveAur(module map[string]interface{}) modules.Aur {
	aur := modules.Aur{}

	if value, ok := module["name"]; ok {
		if name, ok := value.(string); ok {
			aur.Name = []string{name}
		} else if nameValues, ok := value.([]interface{}); ok {
			names := []string{}
			for _, nameValue := range nameValues {
				if name, ok := nameValue.(string); ok {
					names = append(names, name)
				}
			}
			aur.Name = names
		}
	}
	return aur
}

func retrieveSymlink(modulePath string, module map[string]interface{}) modules.Symlink {
	symlink := modules.Symlink{ModulePath: modulePath}

	if value, ok := module["sudo"]; ok {
		if sudo, ok := value.(bool); ok {
			symlink.Sudo = sudo
		}
	}

	if value, ok := module["from"]; ok {
		if from, ok := value.(string); ok {
			symlink.From = from
		}
	}

	if value, ok := module["to"]; ok {
		if to, ok := value.(string); ok {
			symlink.To = to
		}
	}

	return symlink
}

func retrieveGroup(module map[string]interface{}) modules.Group {
	group := modules.Group{}

	if value, ok := module["user"].(string); ok {
		group.User = value
	}
	if value, ok := module["name"].(string); ok {
		group.Name = value
	}

	return group
}

func retrieveScript(module map[string]interface{}) modules.Script {
	script := modules.Script{}

	if value, ok := module["command"].(string); ok {
		script.Command = value
	}

	return script
}

func retrieveService(module map[string]interface{}) modules.Service {
	service := modules.Service{}

	if value, ok := module["name"].(string); ok {
		service.Name = value
	}

	if value, ok := module["enable"].(bool); ok {
		service.Enable = value
	}

	if value, ok := module["start"].(bool); ok {
		service.Start = value
	}
	return service
}

func retrieveDirectory(module map[string]interface{}) modules.Directory {
	directory := modules.Directory{}

	if value, ok := module["sudo"].(bool); ok {
		directory.Sudo = value
	}

	if value, ok := module["path"].(string); ok {
		directory.Path = value
	}

	return directory
}
