package main

import "fmt"

import "log"
import "path/filepath"
import "path"
import "os"
import "gopkg.in/yaml.v2"
import "os/exec"
import "strings"
import "golang.org/x/exp/slices"

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

type Directory struct {
	Path string
}

type Service struct {
	Name   string
	Enable bool
	Start  bool
}

type Group struct {
	User  string
	Group string
}

type Line struct {
	Path string
	Line string
}

type Repository struct {
	Repository  string
	Destination string
}

type Fragment struct {
	Path       string
	Pacman     Pacman
	Aur        Aur
	Symlink    []Symlink
	Script     []Script
	Shell      Shell
	Directory  []Directory
	Service    []Service
	Group      []Group
	Line       []Line
	Repository []Repository
}

type Configuration struct {
	Fragments []Fragment
}

func installed_packages() ([]string, error) {
	bytes, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
}

func (configuration Configuration) execute() error {
	pacman_args := []string{"pacman", "-S"}
	installed, err := installed_packages()
	if err != nil {
		return err
	}

	for _, fragment := range configuration.Fragments {
		fmt.Println(fragment.Path)
		for _, pacman_package := range fragment.Pacman.Name {
			if !slices.Contains(installed, pacman_package) {
				pacman_args = append(pacman_args, pacman_package)
			}
		}
	}

	cmd := exec.Command("sudo", pacman_args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func main() {
	module_paths, err := filepath.Glob("../../.dotfiles/fragments/*")

	configuration := Configuration{}

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, module_path := range module_paths {
		configuration_path := path.Join(module_path, "configuration.yml")

		contents, err := readFile(configuration_path)
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

	fmt.Printf("--- t:\n%+v\n\n", configuration)
	err = configuration.execute()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
