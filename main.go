package main

import "fmt"
import "errors"
import "log"
import "path/filepath"
import "os"
import "gopkg.in/yaml.v2"
import "os/user"
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

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

type Pacman struct {
	Name []string
}

type Aur struct {
	Name []string
}

type Symlink struct {
	Sudo     bool
	From     string
	To       string
	fullFrom string
	fullTo   string
}

type Script struct {
	Command string
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
	User string
	Name string
}

type Fragment struct {
	Path      string
	Pacman    Pacman
	Aur       Aur
	Symlink   []Symlink
	Script    []Script
	Directory []Directory
	Service   []Service
	Group     []Group
}

type Configuration struct {
	Fragments []Fragment
}

func installedPackages() ([]string, error) {
	bytes, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
}

func (configuration Configuration) execute() error {
	err := configuration.executePacman()
	if err != nil {
		return err
	}

	err = configuration.executeAur()
	if err != nil {
		return err
	}

	err = configuration.executeSymlink()
	if err != nil {
		return err
	}

	err = configuration.executeScript()
	if err != nil {
		return err
	}

	err = configuration.executeDirectory()
	if err != nil {
		return err
	}

	err = configuration.executeService()
	if err != nil {
		return err
	}

	err = configuration.executeGroup()
	if err != nil {
		return err
	}

	return nil
}

func (configuration Configuration) executePacman() error {
	missing := []string{}
	installed, err := installedPackages()
	if err != nil {
		return err
	}

	for _, fragment := range configuration.Fragments {
		for _, p := range fragment.Pacman.Name {
			if !slices.Contains(installed, p) {
				missing = append(missing, p)
			}
		}
	}

	if len(missing) > 0 {
		cmd := exec.Command("sudo", append([]string{"pacman", "-S"}, missing...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func (configuration Configuration) executeAur() error {
	missing := []string{}
	installed, err := installedPackages()
	if err != nil {
		return err
	}

	for _, fragment := range configuration.Fragments {
		for _, p := range fragment.Aur.Name {
			if !slices.Contains(installed, p) {
				missing = append(missing, p)
			}
		}
	}

	if len(missing) > 0 {
		cmd := exec.Command("yay", append([]string{"-S"}, missing...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func expandTilde(path string) string {
	usr, _ := user.Current()
	home := usr.HomeDir

	if path == "~" {
		return home
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}

	return path
}

func isSymlinked(from string, to string) bool {
	destination, err := os.Readlink(from)

	if err != nil {
		return false
	}

	return destination == to
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrExist)
}

func (configuration Configuration) executeSymlink() error {
	missing := []Symlink{}

	for _, fragment := range configuration.Fragments {
		for _, symlink := range fragment.Symlink {
			fullFrom := expandTilde(symlink.From)
			fullTo, err := filepath.Abs(filepath.Join(fragment.Path, "files", symlink.To))

			if err != nil {
				return err
			}

			if strings.HasPrefix(symlink.To, "/") {
				fullTo = symlink.To
			}
			symlink.fullFrom = fullFrom
			symlink.fullTo = fullTo

			if !isSymlinked(symlink.fullFrom, symlink.fullTo) {
				missing = append(missing, symlink)
			}
		}
	}

	for _, symlink := range missing {
		if fileExists(symlink.fullFrom) {
			return errors.New(fmt.Sprintf("File \"%s\" already exists"))
		} else {
			fmt.Println("Symlinking " + symlink.fullFrom + " => " + symlink.fullTo)

			cmd := exec.Command("ln", "-s", symlink.fullTo, symlink.fullFrom)
			if symlink.Sudo {
				cmd = exec.Command("sudo", "ln", "-s", symlink.fullTo, symlink.fullFrom)
			}

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (configuration Configuration) executeScript() error {
	for _, fragment := range configuration.Fragments {
		for _, script := range fragment.Script {
			cmd := exec.Command("bash", "-c", script.Command)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (configuration Configuration) executeDirectory() error {
	for _, fragment := range configuration.Fragments {
		for _, directory := range fragment.Directory {
			fullPath := expandTilde(directory.Path)

			isDirectory, err := isDirectory(fullPath)
			if err != nil {
				return err
			}

			if !isDirectory {
				cmd := exec.Command("mkdir", fullPath)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s Service) IsStarted() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 0, nil
		}
	}

	return true, nil
}

func (s Service) IsEnabled() (bool, error) {
	cmd := exec.Command("systemctl", "is-enabled", "--quiet", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 0, nil
		}
	}

	return true, nil
}

func (s Service) StartService() error {
	cmd := exec.Command("sudo", "systemctl", "start", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (s Service) EnableService() error {
	cmd := exec.Command("sudo", "systemctl", "enable", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (configuration Configuration) executeService() error {
	for _, fragment := range configuration.Fragments {
		for _, service := range fragment.Service {
			if service.Enable {
				enabled, err := service.IsEnabled()
				if err != nil {
					return err
				}
				if !enabled {
					fmt.Println("Enabling service \"" + service.Name + "\"")
					err := service.EnableService()
					if err != nil {
						return err
					}
				}
			}

			if service.Start {
				active, err := service.IsStarted()
				if err != nil {
					return err
				}
				if !active {
					fmt.Println("Starting service \"" + service.Name + "\"")
					err := service.StartService()
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (g Group) IsUserInGroup(usr *user.User) (bool, error) {
	groupIDs, err := usr.GroupIds()

	if err != nil {
		return false, err
	}

	for _, groupID := range groupIDs {
		group, err := user.LookupGroupId(groupID)

		if err != nil {
			return false, err
		}

		if group.Name == g.Name {
			return true, nil
		}
	}

	return false, nil
}

func (g Group) AddToUser(usr *user.User) error {
	cmd := exec.Command("sudo", "usermod", "-a", "-G", g.Name, usr.Username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (configuration Configuration) executeGroup() error {
	for _, fragment := range configuration.Fragments {
		for _, group := range fragment.Group {
			usr, err := user.Lookup(group.User)

			if err != nil {
				return err
			}

			userInGroup, err := group.IsUserInGroup(usr)

			if err != nil {
				return err
			}

			if !userInGroup {
				fmt.Println("Adding user \"" + group.User + "\"to group \"" + group.Name + "\"")
				group.AddToUser(usr)
			} else {
			}
		}
	}

	return nil
}

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

	err = configuration.execute()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
