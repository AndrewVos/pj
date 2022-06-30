package cmd

import (
	"github.com/AndrewVos/pj/applyables"
	"github.com/AndrewVos/pj/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

func buildCommand(applyableInfo CommandInfo) {
	value := reflect.Indirect(reflect.ValueOf(applyableInfo.Applyable))
	var command = &cobra.Command{
		Use:   "add-" + applyableInfo.Name + " <module_name>",
		Short: applyableInfo.Description,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data := map[string]interface{}{}

			cmd.Flags().Visit(func(f *pflag.Flag) {
				for i := 0; i < value.NumField(); i++ {
					originalFieldName := value.Type().Field(i).Name
					fieldName := strings.ToLower(originalFieldName)
					field := value.FieldByName(originalFieldName)

					if field.CanInterface() {
						fieldValue := field.Interface()
						kind := field.Kind().String()

						if f.Name == fieldName {
							data[f.Name] = fieldValue

							if _, ok := fieldValue.(string); ok {
								s, _ := cmd.Flags().GetString(fieldName)
								data[f.Name] = s
							} else if _, ok := fieldValue.([]string); ok {
								s, _ := cmd.Flags().GetStringArray(fieldName)
								data[f.Name] = s
							} else if _, ok := fieldValue.(bool); ok {
								s, _ := cmd.Flags().GetBool(fieldName)
								data[f.Name] = s
							} else {
								log.Fatalf("unsupported flag field type %v\n", kind)
							}
						}
					}
				}
			})

			moduleName := args[0]
			modulePath := path.Join("modules", moduleName)
			configurationPath := path.Join(modulePath, "configuration.yml")

			err := os.MkdirAll(modulePath, 0777)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			b, err := yaml.Marshal([]map[string]interface{}{map[string]interface{}{applyableInfo.Name: data}})
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			err = utils.AppendToFile(configurationPath, string(b))
			if err != nil {
				log.Fatalf("error: %v", err)
			}
		},
	}

	for i := 0; i < value.NumField(); i++ {
		originalFieldName := value.Type().Field(i).Name
		fieldName := strings.ToLower(originalFieldName)
		field := value.FieldByName(originalFieldName)

		if field.CanInterface() {
			fieldValue := field.Interface()
			kind := field.Kind().String()

			if _, ok := fieldValue.(string); ok {
				command.Flags().String(fieldName, "", fieldName)
			} else if _, ok := fieldValue.([]string); ok {
				command.Flags().StringArray(fieldName, []string{}, fieldName)
			} else if _, ok := fieldValue.(bool); ok {
				command.Flags().Bool(fieldName, false, fieldName)
			} else {
				log.Fatalf("unsupported flag field type %v\n", kind)
			}
		}
	}

	command.Flags().SortFlags = false

	rootCmd.AddCommand(command)
}

type CommandInfo struct {
	Name        string
	Applyable   applyables.Applyable
	Description string
}

func init() {
	commandInfos := []CommandInfo{
		CommandInfo{Name: "aur", Applyable: applyables.Aur{}, Description: "Add an AUR package"},
		CommandInfo{Name: "brew", Applyable: applyables.Brew{}, Description: "Add a Homebrew package"},
		CommandInfo{Name: "pacman", Applyable: applyables.Pacman{}, Description: "Add a Pacman package"},
		CommandInfo{Name: "directory", Applyable: applyables.Directory{}, Description: "Add a Directory"},
		CommandInfo{Name: "group", Applyable: applyables.Group{}, Description: "Add a Group"},
		CommandInfo{Name: "script", Applyable: applyables.Script{}, Description: "Add a Script"},
		CommandInfo{Name: "service", Applyable: applyables.Service{}, Description: "Add a Service"},
		CommandInfo{Name: "symlink", Applyable: applyables.Symlink{}, Description: "Add a Symlink"},
	}

	for _, i := range commandInfos {
		buildCommand(i)
	}
}
