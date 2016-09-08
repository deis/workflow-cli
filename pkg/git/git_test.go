package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/arschles/assert"
)

func TestRepositoryURL(t *testing.T) {
	t.Parallel()

	actual := RepositoryURL("deis.example.com", "app")
	assert.Equal(t, actual, "ssh://git@deis-builder.example.com:2222/app.git", "url")
	actual = RepositoryURL("deis.10.245.1.3.xip.io:31350", "velcro-underdog")
	assert.Equal(t, actual, "ssh://git@deis-builder.10.245.1.3.xip.io:2222/velcro-underdog.git", "url")
}

func TestGetRemotes(t *testing.T) {
	t.Parallel()

	expected := []remote{
		{"test", "ssh://test.com/test.git"},
		{"example", "ssh://example.com:2222/example.git"},
	}

	actual, err := getRemotes(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
example	ssh://example.com:2222/example.git (fetch)
example	ssh://example.com:2222/example.git (push)
`, nil
	})

	assert.NoErr(t, err)
	assert.Equal(t, actual, expected, "remotes")

	_, err = getRemotes(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `fakeoutput(push)
`, nil
	})

	assert.Err(t, err, ErrInvalidRepositoryList)

	testErr := errors.New("test error")

	_, err = getRemotes(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return "", testErr
	})

	assert.Err(t, err, testErr)
}

func TestRemoteURL(t *testing.T) {
	t.Parallel()

	url, err := RemoteURL(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
example	ssh://example.com:2222/example.git (fetch)
example	ssh://example.com:2222/example.git (push)
`, nil
	}, "test")

	assert.NoErr(t, err)
	assert.Equal(t, url, "ssh://test.com/test.git", "remote url")

	_, err = RemoteURL(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
example	ssh://example.com:2222/example.git (fetch)
example	ssh://example.com:2222/example.git (push)
`, nil
	}, "foo")

	assert.Err(t, err, ErrRemoteNotFound)

	testErr := errors.New("test error")

	_, err = RemoteURL(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return "", testErr
	}, "test")

	assert.Err(t, err, testErr)
}

func TestFindRemoteURL(t *testing.T) {
	t.Parallel()

	url, err := findRemote(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
deis	ssh://git@deis-builder.example.com:2222/test.git (fetch)
deis	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
	}, "deis.example.com")

	assert.NoErr(t, err)
	assert.Equal(t, url, "ssh://git@deis-builder.example.com:2222/test.git", "remote url")

	_, err = findRemote(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
example	ssh://example.com:2222/example.git (fetch)
example	ssh://example.com:2222/example.git (push)
`, nil
	}, "deis.test.com")

	assert.Err(t, err, ErrRemoteNotFound)

	testErr := errors.New("test error")

	_, err = findRemote(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return "", testErr
	}, "deis.test.com")

	assert.Err(t, err, testErr)
}

func TestDetectAppName(t *testing.T) {
	t.Parallel()

	app, err := DetectAppName(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `deis	ssh://git@deis-builder.example.com:2222/test.git (fetch)
deis	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
	}, "deis.example.com")

	assert.NoErr(t, err)
	assert.Equal(t, app, "test", "app")

	app, err = DetectAppName(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return "", errors.New("test error")
	}, "deis.test.com")
	wd, err := os.Getwd()
	assert.NoErr(t, err)
	assert.Equal(t, app, filepath.Base(wd), "app")
}

func TestRemoteNamesFromAppID(t *testing.T) {
	t.Parallel()

	apps, err := remoteNamesFromAppID(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
deis	ssh://git@deis-builder.example.com:2222/test.git (fetch)
deis	ssh://git@deis-builder.example.com:2222/test.git (push)
two	ssh://git@deis-builder.example.com:2222/test.git (fetch)
two	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
	}, "deis.example.com", "test")

	assert.NoErr(t, err)
	assert.Equal(t, apps, []string{"deis", "two"}, "remote url")

	_, err = remoteNamesFromAppID(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return `test	ssh://test.com/test.git (fetch)
test	ssh://test.com/test.git (push)
deis	ssh://git@deis-builder.example.com:2222/test.git (fetch)
deis	ssh://git@deis-builder.example.com:2222/test.git (push)
two	ssh://git@deis-builder.example.com:2222/test.git (fetch)
two	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
	}, "deis.test.com", "other")

	assert.Err(t, err, ErrRemoteNotFound)

	testErr := errors.New("test error")

	_, err = remoteNamesFromAppID(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "-v"}, "args")
		return "", testErr
	}, "deis.test.com", "test")

	assert.Err(t, err, testErr)
}

func TestDeleteRemote(t *testing.T) {
	t.Parallel()

	err := DeleteRemote(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "remove", "test"}, "args")
		return "", nil
	}, "test")
	assert.NoErr(t, err)
}

func TestDeleteAppRemotes(t *testing.T) {
	t.Parallel()

	err := DeleteAppRemotes(func(cmd []string) (string, error) {
		if reflect.DeepEqual(cmd, []string{"remote", "-v"}) {
			return `deis	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
		} else if reflect.DeepEqual(cmd, []string{"remote", "remove", "deis"}) {
			return "", nil
		} else {
			t.Errorf("unexpected command %v", cmd)
			return "", nil
		}
	}, "deis.example.com", "test")

	assert.NoErr(t, err)

	testErr := errors.New("test error")

	err = DeleteAppRemotes(func(cmd []string) (string, error) {
		return "", testErr
	}, "deis.example.com", "test")

	assert.Err(t, testErr, err)

	err = DeleteAppRemotes(func(cmd []string) (string, error) {
		if reflect.DeepEqual(cmd, []string{"remote", "-v"}) {
			return `deis	ssh://git@deis-builder.example.com:2222/test.git (push)
`, nil
		} else if reflect.DeepEqual(cmd, []string{"remote", "remove", "deis"}) {
			return "", testErr
		} else {
			t.Errorf("unexpected command %v", cmd)
			return "", nil
		}
	}, "deis.example.com", "test")

	assert.Err(t, testErr, err)
}

func TestInit(t *testing.T) {
	t.Parallel()

	err := Init(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"init"}, "args")
		return "", nil
	})
	assert.NoErr(t, err)
}

func TestCreateRemote(t *testing.T) {
	t.Parallel()

	err := CreateRemote(func(cmd []string) (string, error) {
		assert.Equal(t, cmd, []string{"remote", "add", "deis", "ssh://git@deis-builder.example.com:2222/testing.git"}, "args")
		return "", nil
	}, "deis.example.com", "deis", "testing")
	assert.NoErr(t, err)
}

func TestGitError(t *testing.T) {
	t.Parallel()

	exitErr := exec.ExitError{Stderr: []byte("fake error")}
	assert.Equal(t, gitError(&exitErr, []string{"fake"}).Error(), `Error when running 'git fake'
fake error`, "error")
}
