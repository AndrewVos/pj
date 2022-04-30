package modules

import "os"
import "os/exec"

type Service struct {
	Name   string
	Enable bool
	Start  bool
}

func (s Service) IsStarted() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 0, nil
		}
	}

	return true, nil
}

func (s Service) IsEnabled() (bool, error) {
	cmd := exec.Command("systemctl", "is-enabled", "--quiet", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() == 0, nil
		}
	}

	return true, nil
}

func (s Service) StartService() error {
	cmd := exec.Command("sudo", "systemctl", "start", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (s Service) EnableService() error {
	cmd := exec.Command("sudo", "systemctl", "enable", s.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
