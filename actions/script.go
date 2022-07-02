package actions

import (
	"os"
	"os/exec"
)

type Script struct {
	Command string `flag:"required"`
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

func (a Script) Completions(fieldName string) ([]string, error) {
	return []string{}, nil
}

func (s Script) Apply(modulePath string) error {
	cmd := exec.Command("bash", "-c", s.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
