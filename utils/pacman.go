package utils

import "strings"
import "os/exec"

var allPacmanPackages []string

func ListInstalledPackages() ([]string, error) {
	bytes, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
}

func ListPacmanPackages() ([]string, error) {
	if len(allPacmanPackages) == 0 {
		bytes, err := exec.Command("pacman", "-Ssq").Output()
		if err != nil {
			return nil, err
		}
		allPacmanPackages = strings.Split(strings.TrimSpace(string(bytes)), "\n")
	}
	return allPacmanPackages, nil
}
