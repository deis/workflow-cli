package cmd

import (
	"fmt"

	"github.com/deis/workflow-cli/pkg/git"
)

const remoteCreationMsg = "Git remote %s successfully created for app %s.\n"
const remoteDeletionMsg = "Git remotes for app %s removed.\n"

// GitRemote creates a git remote for a deis app.
func (d *DeisCmd) GitRemote(appID, remote string, force bool) error {
	s, appID, err := load(d.ConfigFile, appID)

	remoteURL, err := git.RemoteURL(git.DefaultCmd, remote)

	if err != nil {
		//If git remote doesn't exist, create it without issue
		if err == git.ErrRemoteNotFound {
			err := git.CreateRemote(git.DefaultCmd, s.Client.ControllerURL.Host, remote, appID)
			if err != nil {
				return err
			}
			d.Printf(remoteCreationMsg, remote, appID)
			return nil
		}

		return err
	}

	expectedURL := git.RepositoryURL(s.Client.ControllerURL.Host, appID)

	if remoteURL == expectedURL {
		d.Printf("Remote %s already exists and is correctly configured for app %s.\n", remote, appID)
		return nil
	}

	if force {
		d.Printf("Deleting git remote %s.\n", remote)
		err := git.DeleteRemote(git.DefaultCmd, remote)
		if err != nil {
			return err
		}
		err = git.CreateRemote(git.DefaultCmd, s.Client.ControllerURL.Host, remote, appID)
		if err != nil {
			return err
		}
		d.Printf(remoteCreationMsg, remote, appID)
		return nil
	}

	msg := "Remote %s already exists, please run 'deis git:remote -f' to overwrite\n"
	msg += "Existing remote URL: %s\n"
	msg += "When forced, will overwrite with: %s"

	return fmt.Errorf(msg, remote, remoteURL, expectedURL)
}

// GitRemove removes a application git remote from a repository
func (d *DeisCmd) GitRemove(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	err = git.DeleteAppRemotes(git.DefaultCmd, s.Client.ControllerURL.Host, appID)

	if err != nil {
		return err
	}

	d.Printf(remoteDeletionMsg, appID)
	return nil
}
