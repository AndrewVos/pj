package main

import "fmt"

import "log"
import "path/filepath"
import "path"
import "os"
import "gopkg.in/yaml.v2"

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Pacman struct {
	Name []string
}

type Symlink struct {
	From string
	To   string
}

type Script struct {
	Command string
}

type Shell struct {
	Shell string
}

type Fragment struct {
	Pacman  Pacman
	Symlink []Symlink
	Script  []Script
	Shell   Shell
}

func main() {
	module_paths, err := filepath.Glob("../../.dotfiles/fragments/*")

	if err != nil {
		panic(err)
	}

	for _, module_path := range module_paths {
		configuration_path := path.Join(module_path, "configuration.yml")
		fmt.Println(configuration_path)

		if _, err := os.Stat(configuration_path); err == nil {
			contents := readFile(configuration_path)

			fragment := Fragment{}
			err := yaml.Unmarshal([]byte(contents), &fragment)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			fmt.Printf("--- t:\n%+v\n\n", fragment)
		}
	}
}
