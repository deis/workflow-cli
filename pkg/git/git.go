package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ErrRemoteNotFound is returned when the remote cannot be found in git
var ErrRemoteNotFound = errors.New("Could not find remote matching app in 'git remote -v'")

func gitError(err *exec.ExitError, cmd []string) error {
	msg := fmt.Sprintf("Error when running '%s'\n", strings.Join(cmd, " "))
	out := string(err.Stderr)
	if out != "" {
		msg += strings.TrimSpace(out)
	}

	return errors.New(msg)
}

// CreateRemote adds a git remote in the current directory.
func CreateRemote(host, remote, appID string) error {
	cmd := []string{"git", "remote", "add", remote, RemoteURL(host, appID)}
	_, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return gitError(err.(*exec.ExitError), cmd)
	}

	return nil
}

// DeleteAppRemotes removes all git remotes corresponding to an app in the repository.
func DeleteAppRemotes(host, appID string) error {
	names, err := remoteNamesFromAppID(host, appID)

	if err != nil {
		return err
	}

	for _, name := range names {
		if err := DeleteRemote(name); err != nil {
			return err
		}
	}

	return nil
}

// DeleteRemote removes a remote from the repository
func DeleteRemote(name string) error {
	cmd := []string{"git", "remote", "remove", name}
	_, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return gitError(err.(*exec.ExitError), cmd)
	}

	return nil
}

// remoteNamesFromAppID returns the git remote names for an app
func remoteNamesFromAppID(host, appID string) ([]string, error) {
	cmd := []string{"git", "remote", "-v"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()

	if err != nil {
		return []string{}, gitError(err.(*exec.ExitError), cmd)
	}

	remotes := []string{}

lines:
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, RemoteURL(host, appID)) {
			name := strings.Split(line, "\t")[0]
			// git remote -v can show duplicate remotes, so don't add a remote if it already has been added
			for _, remote := range remotes {
				if remote == name {
					continue lines
				}
			}
			remotes = append(remotes, name)
		}
	}

	if len(remotes) == 0 {
		return remotes, ErrRemoteNotFound
	}

	return remotes, nil
}

// DetectAppName detects if there is deis remote in git.
func DetectAppName(host string) (string, error) {
	remote, err := findRemote(host)

	// Don't return an error if remote can't be found, return directory name instead.
	if err != nil {
		dir, err := os.Getwd()
		return strings.ToLower(filepath.Base(dir)), err
	}

	ss := strings.Split(remote, "/")
	return strings.Split(ss[len(ss)-1], ".")[0], nil
}

func findRemote(host string) (string, error) {
	cmd := []string{"git", "remote", "-v"}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", gitError(err.(*exec.ExitError), cmd)
	}

	// Strip off any trailing :port number after the host name.
	host = strings.Split(host, ":")[0]
	builderHost := getBuilderHostname(host)

	for _, line := range strings.Split(string(out), "\n") {
		for _, remote := range strings.Split(line, " ") {
			if strings.Contains(remote, host) || strings.Contains(remote, builderHost) {
				return strings.Split(remote, "\t")[1], nil
			}
		}
	}

	return "", ErrRemoteNotFound
}

// RemoteURL returns the git URL of app.
func RemoteURL(host, appID string) string {
	// Strip off any trailing :port number after the host name.
	host = strings.Split(host, ":")[0]
	return fmt.Sprintf("ssh://git@%s:2222/%s.git", getBuilderHostname(host), appID)
}

// getBuilderHostname derives the builder host name from the controller host name.
func getBuilderHostname(host string) string {
	hostTokens := strings.Split(host, ".")
	hostTokens[0] = fmt.Sprintf("%s-builder", hostTokens[0])
	return strings.Join(hostTokens, ".")
}

// RemoteValue gets the url that a git remote is set to.
func RemoteValue(name string) (string, error) {
	cmd := []string{"git", "remote", "get-url", name}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()

	if err != nil {
		// get the return code of the program and see if it equals not found
		if err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus() == 128 {
			return "", ErrRemoteNotFound
		}
		return "", gitError(err.(*exec.ExitError), cmd)
	}

	return strings.Trim(string(out), "\n"), nil
}
