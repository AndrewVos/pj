package utils

import "strings"
import "os/exec"

func ListInstalledBrews() ([]string, error) {
	bytes, err := exec.Command("brew", "list").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
}
