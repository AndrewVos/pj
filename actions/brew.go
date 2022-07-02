package actions

import (
	"github.com/AndrewVos/pj/utils"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
)

type Brew struct {
	Name []string
}

func init() {
	RegisterAction(Brew{})
}

func (a Brew) Flag() string {
	return "brew"
}

func (a Brew) AddActionDescription() string {
	return "Add a Homebrew package"
}

func (p Brew) Apply(modulePath string) error {
	missing := []string{}

	installed, err := utils.ListInstalledBrews()
	if err != nil {
		return err
	}

	for _, pkg := range p.Name {
		if !slices.Contains(installed, pkg) {
			missing = append(missing, pkg)
		}
	}

	if len(missing) > 0 {
		cmd := exec.Command("brew", append([]string{"install"}, missing...)...)
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
