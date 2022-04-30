package utils

import "os"
import "os/user"
import "strings"
import "path/filepath"
import "errors"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrExist)
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

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
