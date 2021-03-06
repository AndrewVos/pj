package tasks

import (
	"github.com/AndrewVos/pj/utils"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
)

type Aur struct {
	Name []string `flag:"required"`
}

func init() {
	RegisterTask(Aur{})
}

func (a Aur) Flag() string {
	return "aur"
}

func (a Aur) AddTaskDescription() string {
	return "Add an AUR package"
}

func (a Aur) Completions(fieldName string) ([]string, error) {
	return []string{}, nil
}

func (a Aur) Apply(modulePath string) error {
	missing := []string{}
	installed, err := utils.ListInstalledPackages()

	if err != nil {
		return err
	}

	for _, pkg := range a.Name {
		if !slices.Contains(installed, pkg) {
			missing = append(missing, pkg)
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
