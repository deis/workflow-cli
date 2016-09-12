// +build linux darwin

package settings

import (
	"os"
	"testing"

	"github.com/arschles/assert"
)

// TestFindHome ensures the correct home directory is returned by FindHome().
func TestFindHome(t *testing.T) {
	expectedHomeDir := "/d/e/f"
	os.Setenv("HOME", expectedHomeDir)

	assert.Equal(t, FindHome(), expectedHomeDir, "output")
}

// TestSetHome ensures the correct env vars are set when SetHome() is called.
func TestSetHome(t *testing.T) {
	homeDir := "/a/b/c"
	SetHome(homeDir)

	assert.Equal(t, os.Getenv("HOME"), homeDir, "output")
}
