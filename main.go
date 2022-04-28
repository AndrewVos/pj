package main

import "fmt"

import "log"
import "path/filepath"
import "path"
import "os"
import "gopkg.in/yaml.v2"

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Pacman struct {
	Name []string
}

type Aur struct {
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
	Aur     Aur
	Symlink []Symlink
	Script  []Script
	Shell   Shell
}

func main() {
	module_paths, err := filepath.Glob("../../.dotfiles/fragments/*")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, module_path := range module_paths {
		configuration_path := path.Join(module_path, "configuration.yml")
		fmt.Println(configuration_path)

		contents, err := readFile(configuration_path)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fragment := Fragment{}
		err = yaml.Unmarshal([]byte(contents), &fragment)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("--- t:\n%+v\n\n", fragment)
	}
}
