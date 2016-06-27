package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/users"
	"github.com/deis/workflow-cli/settings"
)

// UsersList lists users registered with the controller.
func UsersList(results int) error {
	s, err := settings.Load()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	users, count, err := users.List(s.Client, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== Users%s", limitCount(len(users), count))

	for _, user := range users {
		fmt.Println(user.Username)
	}
	return nil
}
