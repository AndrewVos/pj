package cmd

import (
	"github.com/AndrewVos/pj/actions"
	"github.com/AndrewVos/pj/modules"
	"github.com/AndrewVos/pj/utils"
	"github.com/fatih/structtag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

func buildCommand(action actions.Action) *cobra.Command {
	value := reflect.Indirect(reflect.ValueOf(action))

	var command = &cobra.Command{
		Use:   action.Flag() + " <module_name>",
		Short: action.AddActionDescription(),
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

			b, err := yaml.Marshal([]map[string]interface{}{map[string]interface{}{action.Flag(): data}})
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			err = utils.AppendToFile(configurationPath, string(b))
			if err != nil {
				log.Fatalf("error: %v", err)
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}
			moduleNames, err := getModuleNames()
			if err != nil {
				panic(err)
			}
			return moduleNames, cobra.ShellCompDirectiveNoFileComp
		},
	}

	for i := 0; i < value.NumField(); i++ {
		originalFieldName := value.Type().Field(i).Name
		tag := value.Type().Field(i).Tag
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

			tags, err := structtag.Parse(string(tag))
			if err != nil {
				panic(err)
			}

			flagTag, err := tags.Get("flag")
			if err == nil {
				if flagTag.Name == "required" {
					command.MarkFlagRequired(fieldName)
				}
			}

			completionTag, err := tags.Get("completion")
			if err == nil {
				command.RegisterFlagCompletionFunc(fieldName, completionFunc(action, completionTag.Name))
			}
		}
	}

	command.Flags().SortFlags = false

	return command
}

func completionFunc(action actions.Action, name string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		completions, err := action.Completions(name)
		if err != nil {
			panic(err)
		}
		return completions, cobra.ShellCompDirectiveDefault
	}
}

func getModuleNames() ([]string, error) {
	var names []string

	modules, err := modules.LoadModules()
	if err != nil {
		return names, err
	}

	for _, m := range modules {
		names = append(names, m.Name)
	}
	return names, err
}

func init() {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add an action",
	}

	for _, action := range actions.All {
		addCmd.AddCommand(buildCommand(action))
	}
	rootCmd.AddCommand(addCmd)
}
