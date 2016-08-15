package cmd

import (
	"github.com/deis/controller-sdk-go/users"
	"github.com/deis/workflow-cli/settings"
)

// UsersList lists users registered with the controller.
func (d DeisCmd) UsersList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	users, count, err := users.List(s.Client, results)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== Users%s", limitCount(len(users), count))

	for _, user := range users {
		d.Println(user.Username)
	}
	return nil
}
