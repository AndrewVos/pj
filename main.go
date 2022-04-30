package main

import "log"
import "path/filepath"
import "os"
import "gopkg.in/yaml.v2"
import "github.com/AndrewVos/pj/utils"

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	module_paths, err := filepath.Glob(filepath.Join(cwd, "modules", "*"))

	configuration := Configuration{}

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, module_path := range module_paths {
		configuration_path := filepath.Join(module_path, "configuration.yml")

		contents, err := utils.ReadFile(configuration_path)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fragment := Fragment{Path: module_path}
		err = yaml.Unmarshal([]byte(contents), &fragment)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		configuration.Fragments = append(configuration.Fragments, fragment)
	}

	err = configuration.execute()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
