package actions

import (
	"os"
	"os/exec"
)

type Script struct {
	Command string
}

func init() {
	RegisterAction(Script{})
}

func (a Script) Flag() string {
	return "script"
}

func (a Script) AddActionDescription() string {
	return "Add a Script"
}

func (s Script) Apply() error {
	cmd := exec.Command("bash", "-c", s.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
