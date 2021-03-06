package tasks

import (
	"github.com/AndrewVos/pj/utils"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
)

type Pacman struct {
	Name []string `flag:"required" completion:"pacman-packages"`
}

func init() {
	RegisterTask(Pacman{})
}

func (a Pacman) Flag() string {
	return "pacman"
}

func (a Pacman) AddTaskDescription() string {
	return "Add a Pacman package"
}

func (a Pacman) Completions(name string) ([]string, error) {
	if name == "pacman-packages" {
		var results []string
		results, err := utils.ListPacmanPackages()
		return results, err
	}
	return []string{}, nil
}

func (p Pacman) Apply(modulePath string) error {
	missing := []string{}

	installed, err := utils.ListInstalledPackages()
	if err != nil {
		return err
	}

	for _, pkg := range p.Name {
		if !slices.Contains(installed, pkg) {
			missing = append(missing, pkg)
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
