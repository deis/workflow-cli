package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/perms"
	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/settings"
)

// PermsList prints which users have permissions.
func PermsList(appID string, admin bool, results int) error {
	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	var users []string
	var count int

	if admin {
		if results == defaultLimit {
			results = c.ResponseLimit
		}
		users, count, err = perms.ListAdmins(c, results)
	} else {
		users, err = perms.List(c, appID)
	}

	if checkAPICompatibility(c, err) != nil {
		return err
	}

	if admin {
		fmt.Printf("=== Administrators%s", limitCount(len(users), count))
	} else {
		fmt.Printf("=== %s's Users\n", appID)
	}

	for _, user := range users {
		fmt.Println(user)
	}

	return nil
}

// PermCreate adds a user to an app or makes them an administrator.
func PermCreate(appID string, username string, admin bool) error {

	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	if admin {
		fmt.Printf("Adding %s to system administrators... ", username)
		err = perms.NewAdmin(c, username)
	} else {
		fmt.Printf("Adding %s to %s collaborators... ", username, appID)
		err = perms.New(c, appID, username)
	}

	if checkAPICompatibility(c, err) != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

// PermDelete removes a user from an app or revokes admin privileges.
func PermDelete(appID string, username string, admin bool) error {

	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	if admin {
		fmt.Printf("Removing %s from system administrators... ", username)
		err = perms.DeleteAdmin(c, username)
	} else {
		fmt.Printf("Removing %s from %s collaborators... ", username, appID)
		err = perms.Delete(c, appID, username)
	}

	if checkAPICompatibility(c, err) != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

func permsLoad(appID string, admin bool) (*deis.Client, string, error) {
	c, err := settings.Load()

	if err != nil {
		return nil, "", err
	}

	if !admin && appID == "" {
		appID, err = git.DetectAppName(c.ControllerURL.Host)

		if err != nil {
			return nil, "", err
		}
	}

	return c, appID, err
}
