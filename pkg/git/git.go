package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ErrRemoteNotFound is returned when the remote cannot be found in git
var ErrRemoteNotFound = errors.New("Could not find remote matching app in 'git remote -v'")

// CreateRemote adds a git remote in the current directory.
func CreateRemote(host, remote, appID string) error {
	cmd := exec.Command("git", "remote", "add", remote, RemoteURL(host, appID))
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	output, _ := ioutil.ReadAll(stderr)
	fmt.Print(string(output))

	return cmd.Wait()
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
	_, err := exec.Command("git", "remote", "remove", name).Output()
	return err
}

// remoteNamesFromAppID returns the git remote names for an app
func remoteNamesFromAppID(host, appID string) ([]string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return []string{}, err
	}

	cmd := string(out)
	remotes := []string{}

lines:
	for _, line := range strings.Split(cmd, "\n") {
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
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	// Strip off any trailing :port number after the host name.
	host = strings.Split(host, ":")[0]
	builderHost := getBuilderHostname(host)

	for _, line := range strings.Split(cmd, "\n") {
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
	out, err := exec.Command("git", "remote", "get-url", name).Output()

	if err != nil {
		// get the return code of the program and see if it equals not found
		if err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus() == 128 {
			return "", ErrRemoteNotFound
		}
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}
