package utils

import "strings"
import "os/exec"

func ListInstalledPackages() ([]string,error) {
	bytes, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
}
