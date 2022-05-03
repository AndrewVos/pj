package modules

import (
	"fmt"
	"os"
	"os/exec"
)

type Service struct {
	Name   string
	Enable bool
	Start  bool
}

func (s Service) Apply() error {
	if s.Enable {
		enabled, err := s.IsEnabled()
		if err != nil {
			return err
		}
		if !enabled {
			fmt.Println("Enabling service \"" + s.Name + "\"")
			err := s.EnableService()
			if err != nil {
				return err
			}
		}
	}

	if s.Start {
		active, err := s.IsStarted()
		if err != nil {
			return err
		}
		if !active {
			fmt.Println("Starting service \"" + s.Name + "\"")
			err := s.StartService()
			if err != nil {
				return err
			}
		}
	}

	return nil
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
