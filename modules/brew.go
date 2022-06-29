package modules

import (
	"github.com/AndrewVos/pj/utils"
	"golang.org/x/exp/slices"
	"os"
	"os/exec"
)

type Brew struct {
	Name []string
}

func (p Brew) Apply() error {
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
