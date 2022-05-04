package modules

import (
	"errors"
	"fmt"
	"github.com/AndrewVos/pj/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Symlink struct {
	ModulePath string
	Sudo       bool
	From       string
	To         string
}

func (s Symlink) Apply() error {
	fullFrom := utils.ExpandTilde(s.From)
	fullTo, err := filepath.Abs(filepath.Join(s.ModulePath, "files", s.To))

	if err != nil {
		return err
	}

	if strings.HasPrefix(s.To, "/") {
		fullTo = s.To
	}

	if !utils.IsSymlinked(fullFrom, fullTo) {
		if utils.FileExists(fullFrom) {
			return errors.New(fmt.Sprintf("File \"%s\" already exists"))
		} else {
			fmt.Println("Symlinking " + fullFrom + " => " + fullTo)

			cmd := exec.Command("ln", "-s", fullTo, fullFrom)
			if s.Sudo {
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

	return nil
}
