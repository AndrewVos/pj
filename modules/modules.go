package modules

import (
	"github.com/AndrewVos/pj/tasks"
	"github.com/AndrewVos/pj/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
)

type Module struct {
	Name string `yaml:"-"`
	Path string `yaml:"-"`

	Tasks []map[string]map[string]interface{}
}

func (m Module) Save() error {
	err := os.MkdirAll(m.Path, 0777)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return utils.WriteFile(path.Join(m.Path, "configuration.yml"), string(b))
}

func (m Module) AllTasks() ([]tasks.Task, error) {
	var a []tasks.Task
	for _, data := range m.Tasks {
		task, err := decodeTask(m.Path, data)
		if err != nil {
			return a, err
		}
		a = append(a, task)
	}
	return a, nil
}

func AppendTask(moduleName string, name string, data map[string]interface{}) error {
	modules, err := LoadModules()
	if err != nil {
		return err
	}

	for _, m := range modules {
		if m.Name == moduleName {
			err := m.AppendTask(name, data)
			return err
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	m := Module{Name: moduleName, Path: filepath.Join(cwd, "modules", moduleName)}
	err = m.AppendTask(name, data)
	return err
}

func (m Module) AppendTask(name string, data map[string]interface{}) error {
	m.Tasks = append(m.Tasks, map[string]map[string]interface{}{name: data})
	return m.Save()
}

func (m Module) Apply() error {
	tasks, err := m.AllTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		err := task.Apply(m.Path)
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

func LoadModule(modulePath string) (Module, error) {
	m := Module{Name: path.Base(modulePath), Path: modulePath}

	configurationPath := filepath.Join(modulePath, "configuration.yml")
	if utils.FileExists(configurationPath) {
		contents, err := utils.ReadFile(configurationPath)
		if err != nil {
			return m, err
		}

		err = yaml.Unmarshal([]byte(contents), &m)
		if err != nil {
			return m, err
		}
	}

	return m, nil
}

func LoadModules() ([]Module, error) {
	var modules []Module

	modulePaths, err := modulePaths()
	if err != nil {
		return modules, err
	}

	for _, modulePath := range modulePaths {
		m, err := LoadModule(modulePath)
		if err != nil {
			return modules, err
		}

		modules = append(modules, m)
	}

	return modules, nil
}

func findTaskFor(topLevelModule map[string]map[string]interface{}) (tasks.Task, map[string]interface{}) {
	for _, task := range tasks.All {
		if data, ok := topLevelModule[task.Flag()]; ok {
			return task, data
		}
	}
	panic("task not found")
}

func decodeTask(modulePath string, topLevelModule map[string]map[string]interface{}) (tasks.Task, error) {
	task, data := findTaskFor(topLevelModule)

	if _, ok := task.(tasks.Pacman); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	} else if _, ok := task.(tasks.Aur); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	} else if _, ok := task.(tasks.Brew); ok {
		if name, ok := data["name"].(string); ok {
			data["name"] = []string{name}
		}
	}

	err := mapstructure.Decode(data, &task)
	return task, err
}
