package cmd

import (
	"fmt"

	"github.com/deis/workflow-cli/pkg/git"
)

const remoteCreationMsg = "Git remote %s successfully created for app %s.\n"
const remoteDeletionMsg = "Git remotes for app %s removed.\n"

// GitRemote creates a git remote for a deis app.
func GitRemote(appID string, remote string, force bool) error {
	s, appID, err := load(appID)

	remoteURL, err := git.RemoteValue(remote)

	if err != nil {
		//If git remote doesn't exist, create it without issue
		if err == git.ErrRemoteNotFound {
			git.CreateRemote(s.Client.ControllerURL.Host, remote, appID)
			fmt.Printf(remoteCreationMsg, remote, appID)
			return nil
		}

		return err
	}

	expectedURL := git.RemoteURL(s.Client.ControllerURL.Host, appID)

	if remoteURL == expectedURL {
		fmt.Printf("Remote %s already exists and is correctly configured for app %s.\n", remote, appID)
		return nil
	}

	if force {
		fmt.Printf("Deleting git remote %s.\n", remote)
		git.DeleteRemote(remote)
		git.CreateRemote(s.Client.ControllerURL.Host, remote, appID)
		fmt.Printf(remoteCreationMsg, remote, appID)
		return nil
	}

	msg := "Remote %s already exists, please run 'deis git:remote -f' to overwrite\n"
	msg += "Existing remote URL: %s\n"
	msg += "When forced, will overwrite with: %s"

	return fmt.Errorf(msg, remote, remoteURL, expectedURL)
}

// GitRemove removes a application git remote from a repository
func GitRemove(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	err = git.DeleteAppRemotes(s.Client.ControllerURL.Host, appID)

	if err != nil {
		return err
	}

	fmt.Printf(remoteDeletionMsg, appID)
	return nil
}
