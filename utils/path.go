package utils

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}
}

func ExpandTilde(path string) string {
	usr, _ := user.Current()
	home := usr.HomeDir

	if path == "~" {
		return home
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}

	return path
}

func IsSymlinked(from string, to string) bool {
	destination, err := os.Readlink(from)

	if err != nil {
		return false
	}

	return destination == to
}

func DirectoryExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if errors.Is(err, os.ErrExist) {
		return fileInfo.IsDir(), nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return fileInfo.IsDir(), nil
	}

	return false, err
}
