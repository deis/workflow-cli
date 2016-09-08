package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// ErrRemoteNotFound is returned when the remote cannot be found in git
	ErrRemoteNotFound = errors.New("Could not find remote matching app in 'git remote -v'")
	// ErrInvalidRepositoryList is an error returned if git returns unparsible output
	ErrInvalidRepositoryList = errors.New("Invalid output in 'git remote -v'")
)

// Cmd is a method the exeutes the given git command and returns the output or the error.
type Cmd func(cmd []string) (string, error)

// remote defines a git remote's name and its url.
type remote struct {
	Name string
	URL  string
}

// DefaultCmd is an implementation of Cmd that calls git.
func DefaultCmd(cmd []string) (string, error) {
	out, err := exec.Command("git", cmd...).Output()
	if err != nil {
		return string(out), gitError(err.(*exec.ExitError), cmd)
	}

	return string(out), nil
}

func gitError(err *exec.ExitError, cmd []string) error {
	msg := fmt.Sprintf("Error when running 'git %s'\n", strings.Join(cmd, " "))
	out := string(err.Stderr)
	if out != "" {
		msg += strings.TrimSpace(out)
	}

	return errors.New(msg)
}

// CreateRemote adds a git remote in the current directory.
func CreateRemote(cmd Cmd, host, name, appID string) error {
	_, err := cmd([]string{"remote", "add", name, RepositoryURL(host, appID)})
	return err
}

// Init creates a new git repository in the local directory.
func Init(cmd Cmd) error {
	_, err := cmd([]string{"init"})
	return err
}

// DeleteAppRemotes removes all git remotes corresponding to an app in the repository.
func DeleteAppRemotes(cmd Cmd, host, appID string) error {
	names, err := remoteNamesFromAppID(cmd, host, appID)

	if err != nil {
		return err
	}

	for _, name := range names {
		if err := DeleteRemote(cmd, name); err != nil {
			return err
		}
	}

	return nil
}

// DeleteRemote removes a remote from the repository
func DeleteRemote(cmd Cmd, name string) error {
	_, err := cmd([]string{"remote", "remove", name})
	return err
}

// remoteNamesFromAppID returns the git remote names for an app
func remoteNamesFromAppID(cmd Cmd, host, appID string) ([]string, error) {
	remotes, err := getRemotes(cmd)
	if err != nil {
		return nil, err
	}

	var matchedRemotes []string

	for _, r := range remotes {
		if r.URL == RepositoryURL(host, appID) {
			matchedRemotes = append(matchedRemotes, r.Name)
		}
	}

	if len(matchedRemotes) == 0 {
		return nil, ErrRemoteNotFound
	}

	return matchedRemotes, nil
}

// DetectAppName detects if there is deis remote in git.
func DetectAppName(cmd Cmd, host string) (string, error) {
	remote, err := findRemote(cmd, host)

	// Don't return an error if remote can't be found, return directory name instead.
	if err != nil {
		dir, err := os.Getwd()
		return strings.ToLower(filepath.Base(dir)), err
	}

	ss := strings.Split(remote, "/")
	return strings.Split(ss[len(ss)-1], ".")[0], nil
}

// findRemote finds a remote name the uses a workflow git repository.
func findRemote(cmd Cmd, host string) (string, error) {
	remotes, err := getRemotes(cmd)
	if err != nil {
		return "", err
	}

	// strip port from controller url and use it to find builder hostname
	builderHost := getBuilderHostname(strings.Split(host, ":")[0])

	// search for builder hostname in remote url
	for _, r := range remotes {
		if strings.Contains(r.URL, builderHost) {
			return r.URL, nil
		}
	}

	return "", ErrRemoteNotFound
}

// RepositoryURL returns the git repository of an app.
func RepositoryURL(host, appID string) string {
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

// RemoteURL retrives the url that a git remote is set to.
func RemoteURL(cmd Cmd, name string) (string, error) {
	remotes, err := getRemotes(cmd)
	if err != nil {
		return "", err
	}

	for _, r := range remotes {
		if r.Name == name {
			return r.URL, nil
		}
	}

	return "", ErrRemoteNotFound
}

// getRemotes retrives all the git remotes from a repository
func getRemotes(cmd Cmd) ([]remote, error) {
	out, err := cmd([]string{"remote", "-v"})
	if err != nil {
		return nil, err
	}

	var remotes []remote

	for _, line := range strings.Split(out, "\n") {
		// git remote -v contains both push and fetch remotes.
		// They're generally identical, and deis only cares about push.
		if strings.HasSuffix(line, "(push)") {
			parts := strings.Split(line, "\t")
			if len(parts) < 2 {
				return remotes, ErrInvalidRepositoryList
			}

			remotes = append(remotes, remote{Name: parts[0], URL: strings.Split(parts[1], " ")[0]})
		}
	}

	return remotes, nil
}
