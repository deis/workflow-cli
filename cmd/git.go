package cmd

import (
	"github.com/deis/workflow-cli/pkg/git"
)

// GitRemote creates a git remote for a deis app.
func GitRemote(appID, remote string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	return git.CreateRemote(s.Client.ControllerURL.Host, remote, appID)
}
