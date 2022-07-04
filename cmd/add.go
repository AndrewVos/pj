package cmd

import (
	"github.com/AndrewVos/pj/modules"
	"github.com/AndrewVos/pj/tasks"
	"github.com/fatih/structtag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"reflect"
	"strings"
)

func buildCommand(task tasks.Task) *cobra.Command {
	value := reflect.Indirect(reflect.ValueOf(task))

	var command = &cobra.Command{
		Use:   task.Flag() + " <module_name>",
		Short: task.AddTaskDescription(),
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
			err := modules.AppendTask(moduleName, task.Flag(), data)
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
				command.RegisterFlagCompletionFunc(fieldName, completionFunc(task, completionTag.Name))
			}
		}
	}

	command.Flags().SortFlags = false

	return command
}

func completionFunc(task tasks.Task, name string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		completions, err := task.Completions(name)
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
		Short: "Add a task",
	}

	for _, task := range tasks.All {
		addCmd.AddCommand(buildCommand(task))
	}
	rootCmd.AddCommand(addCmd)
}
