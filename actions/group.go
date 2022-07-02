package actions

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

type Group struct {
	User string `flag:"required"`
	Name string `flag:"required"`
}

func init() {
	RegisterAction(Group{})
}

func (a Group) Flag() string {
	return "group"
}

func (a Group) AddActionDescription() string {
	return "Add a Group"
}

func (a Group) Completions(fieldName string) ([]string, error) {
	return []string{}, nil
}

func (g Group) Apply(modulePath string) error {
	usr, err := user.Lookup(g.User)

	if err != nil {
		return err
	}

	userInGroup, err := g.IsUserInGroup(usr)

	if err != nil {
		return err
	}

	if !userInGroup {
		fmt.Println("Adding user \"" + g.User + "\"to group \"" + g.Name + "\"")
		return g.AddToUser(usr)
	}

	return nil
}

func (g Group) IsUserInGroup(usr *user.User) (bool, error) {
	groupIDs, err := usr.GroupIds()

	if err != nil {
		return false, err
	}

	for _, groupID := range groupIDs {
		group, err := user.LookupGroupId(groupID)

		if err != nil {
			return false, err
		}

		if group.Name == g.Name {
			return true, nil
		}
	}

	return false, nil
}

func (g Group) AddToUser(usr *user.User) error {
	cmd := exec.Command("sudo", "usermod", "-a", "-G", g.Name, usr.Username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
