package modules

import "os"
import "os/user"
import "os/exec"

type Group struct {
	User string
	Name string
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
