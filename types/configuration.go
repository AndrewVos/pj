package types

import "errors"
import "fmt"
import "github.com/AndrewVos/pj/modules"
import "github.com/AndrewVos/pj/utils"
import "golang.org/x/exp/slices"
import "os"
import "os/exec"
import "os/user"
import "path/filepath"
import "strings"

type Configuration struct {
	Fragments []Fragment
}

type Fragment struct {
	Path      string
	Pacman    modules.Pacman
	Aur       modules.Aur
	Symlink   []modules.Symlink
	Script    []modules.Script
	Directory []modules.Directory
	Service   []modules.Service
	Group     []modules.Group
}

func (configuration Configuration) Apply() error {
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
	installed, err := utils.ListInstalledPackages()
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
	installed, err := utils.ListInstalledPackages()
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

func (configuration Configuration) executeSymlink() error {
	for _, fragment := range configuration.Fragments {
		for _, symlink := range fragment.Symlink {
			fullFrom := utils.ExpandTilde(symlink.From)
			fullTo, err := filepath.Abs(filepath.Join(fragment.Path, "files", symlink.To))

			if err != nil {
				return err
			}

			if strings.HasPrefix(symlink.To, "/") {
				fullTo = symlink.To
			}

			if !utils.IsSymlinked(fullFrom, fullTo) {
				if utils.FileExists(fullFrom) {
					return errors.New(fmt.Sprintf("File \"%s\" already exists"))
				} else {
					fmt.Println("Symlinking " + fullFrom + " => " + fullTo)

					cmd := exec.Command("ln", "-s", fullTo, fullFrom)
					if symlink.Sudo {
						cmd = exec.Command("sudo", "ln", "-s", fullTo, fullFrom)
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
			fullPath := utils.ExpandTilde(directory.Path)

			isDirectory, err := utils.IsDirectory(fullPath)
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
