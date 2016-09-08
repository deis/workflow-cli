package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/settings"
)

func TestLoad(t *testing.T) {
	name, err := ioutil.TempDir("", "deis-cli-unit-test-load")
	assert.NoErr(t, err)
	defer os.RemoveAll(name)

	filename := filepath.Join(name, "test.json")
	host := "deis.example.com"
	client, err := deis.New(false, host, "")
	assert.NoErr(t, err)

	config := settings.Settings{
		Username: "test",
		Client:   client,
	}

	filename, err = config.Save(filename)

	_, appID, err := load(filename, "test")
	assert.NoErr(t, err)
	assert.Equal(t, appID, "test", "app")

	assert.NoErr(t, os.Chdir(name))

	_, appID, err = load(filename, "")
	assert.NoErr(t, err)
	assert.Equal(t, appID, filepath.Base(name), "app")

	assert.NoErr(t, git.Init(git.DefaultCmd))
	assert.NoErr(t, git.CreateRemote(git.DefaultCmd, host, "deis", "testing"))

	_, appID, err = load(filename, "")
	assert.NoErr(t, err)
	assert.Equal(t, appID, "testing", "app")
}

func TestDrinkOfChoice(t *testing.T) {
	os.Setenv("DEIS_DRINK_OF_CHOICE", "test")
	assert.Equal(t, drinkOfChoice(), "test", "output")
	os.Unsetenv("DEIS_DRINK_OF_CHOICE")
	assert.Equal(t, drinkOfChoice(), "coffee", "output")
}

func TestLimitsCount(t *testing.T) {
	t.Parallel()
	assert.Equal(t, limitCount(1, 1), "\n", "output")
	assert.Equal(t, limitCount(1, 2), " (1 of 2)\n", "output")
}

func TestAPICompatibility(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	cmdr := DeisCmd{WErr: &b, ConfigFile: ""}
	client := deis.Client{ControllerAPIVersion: "v1.0"}

	err := cmdr.checkAPICompatibility(&client, deis.ErrAPIMismatch)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: 2.3
!    Server version: v1.0
`, "output")

	// After being warned once, the warning should not be printed again.
	b.Reset()
	err = cmdr.checkAPICompatibility(&client, deis.ErrAPIMismatch)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "", "output")

	b.Reset()
	err = cmdr.checkAPICompatibility(&client, deis.ErrConflict)
	assert.Err(t, deis.ErrConflict, err)
	assert.Equal(t, b.String(), "", "output")
}
