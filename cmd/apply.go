package cmd

import (
	"fmt"
	"github.com/AndrewVos/pj/modules"
	"github.com/spf13/cobra"
)

var Verbose bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apply()
	},
}

func init() {
	applyCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(applyCmd)
}

func apply() error {
	modules, err := modules.LoadModules()
	if err != nil {
		return err
	}

	for _, m := range modules {
		if Verbose {
			fmt.Printf("Applying module %s...\n", m.Name)
		}
		m.Apply()
	}

	return err
}
