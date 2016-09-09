// +build linux darwin

package settings

import (
	"os"
	"testing"

	"github.com/arschles/assert"
)

// TestFindHome ensures the correct home directory is returned by FindHome().
func TestFindHome(t *testing.T) {
	previousHome := os.Getenv("HOME")

	expectedHomeDir := "/d/e/f"
	os.Setenv("HOME", expectedHomeDir)

	assert.Equal(t, FindHome(), expectedHomeDir, "output")

	// Reset HOME.
	os.Setenv("HOME", previousHome)
}

// TestSetHome ensures the correct env vars are set when SetHome() is called.
func TestSetHome(t *testing.T) {
	previousHome := os.Getenv("HOME")

	homeDir := "/a/b/c"
	SetHome(homeDir)

	assert.Equal(t, os.Getenv("HOME"), homeDir, "output")

	// Reset HOME.
	os.Setenv("HOME", previousHome)
}
