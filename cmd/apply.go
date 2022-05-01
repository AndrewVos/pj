package cmd

import (
	"fmt"
	"github.com/AndrewVos/pj/types"
	"github.com/AndrewVos/pj/utils"
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

func apply() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	module_paths, err := filepath.Glob(filepath.Join(cwd, "modules", "*"))

	configuration := types.Configuration{}

	if err != nil {
		return err
	}

	for _, module_path := range module_paths {
		configuration_path := filepath.Join(module_path, "configuration.yml")

		contents, err := utils.ReadFile(configuration_path)
		if err != nil {
			return err
		}

		if Verbose {
			moduleName := path.Base(path.Dir(configuration_path))
			fmt.Printf("Applying module %s...\n", moduleName)
		}

		fragment := types.Fragment{Path: module_path}
		err = yaml.Unmarshal([]byte(contents), &fragment)
		if err != nil {
			return err
		}

		configuration.Fragments = append(configuration.Fragments, fragment)
	}

	err = configuration.Apply()
	return err
}

func init() {
	applyCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(applyCmd)
}
