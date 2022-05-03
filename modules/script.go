package modules

import (
	"os"
	"os/exec"
)

type Script struct {
	Command string
}

func (s Script) Apply() error {
	cmd := exec.Command("bash", "-c", s.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
