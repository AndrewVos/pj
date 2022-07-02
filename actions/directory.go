package actions

import (
	"github.com/AndrewVos/pj/utils"
	"os"
	"os/exec"
)

type Directory struct {
	Sudo bool
	Path string `flag:"required"`
}

func init() {
	RegisterAction(Directory{})
}

func (a Directory) Flag() string {
	return "directory"
}

func (a Directory) AddActionDescription() string {
	return "Add a Directory"
}

func (d Directory) Apply(modulePath string) error {
	fullPath := utils.ExpandTilde(d.Path)

	isDirectory, err := utils.DirectoryExists(fullPath)
	if err != nil {
		return err
	}

	if !isDirectory {
		cmd := exec.Command("mkdir", fullPath)
		if d.Sudo {
			cmd = exec.Command("sudo", "mkdir", fullPath)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			return err
		}
	}

	return nil
}
