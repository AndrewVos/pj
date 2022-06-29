package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var moduleName string

var createModuleCmd = &cobra.Command{
	Use:   "create-module <name>",
	Short: "Create a module",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = createDirectory(filepath.Join(cwd, "modules"))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = createDirectory(filepath.Join(cwd, "modules", name))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = createDirectory(filepath.Join(cwd, "modules", name, "files"))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = createConfigurationFile(filepath.Join(cwd, "modules", name, "configuration.yml"))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createModuleCmd)
}

func createDirectory(path string) error {
	fmt.Printf("Creating %s...\n", path)
	return os.MkdirAll(path, 0777)
}

func createConfigurationFile(path string) error {
	fmt.Printf("Creating %s...\n", path)

	content := `### Pacman packages

# Install a single package:
# - pacman:
#     name: imv

# Install multiple packages:

# - pacman:
#     name:
#       - mpv
#       - redis

# Arch User Repository packages

# Install a single package:

# - aur:
#     name: enpass-bin

# Install multiple packages:

# - aur:
#     name:
#       - enpass-bin
#       - google-chrome
#       - spotify
#       - colorpicker
#       - slack-desktop

# Brew packages

# Install a single package:

# - brew:
#     name: postgresql

# Install multiple packages:

# - brew:
#     name:
#       - postgresql
#       - chrome

# Directory

# Create a directory:

# - directory:
#     path: "~/.my-directory"

# Create a directory with sudo:

# - directory:
#     sudo: true
#     path: "/etc/blah"

# Group

# Add a user to a group:

# - group:
#     user: vos
#     name: power

# Script

# Run a script:

# - script:
#     command: '[[ "$SHELL" = /usr/bin/fish ]] || sudo usermod --shell /usr/bin/fish "$USER"'

# Service

# Start a service:

# - service:
#     name: "sshd"
#     start: true

# Start and enable a service:

# - service:
#     name: "sshd"
#     enable: true
#     start: true

# Symlink

# Symlink a file or directory:

# - symlink:
#     from: "~/.ssh"
#     to: ".ssh"

# Symlink a file or directory with sudo:

# - symlink:
#     sudo: true
#     from: "/etc/blah"
#     to: "blah"
`

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	return err
}
