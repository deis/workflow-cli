// +build windows

package settings

import (
	"os"
	"testing"

	"github.com/arschles/assert"
)

// TestFindHome ensures the correct home directory is returned by FindHome().
func TestFindHome(t *testing.T) {
	homedrive := "C:"
	homepath := "/a/b/c"
	os.Setenv("HOMEDRIVE", homedrive)
	os.Setenv("HOMEPATH", homepath)
	assert.Equal(t, FindHome(), homedrive+homepath, "output")
}

// TestSetHome ensures the correct env vars are set when SetHome() is called.
func TestSetHome(t *testing.T) {
	homeDrive := "D:"
	homePath := "/e/f/g"
	SetHome(homeDrive + homePath)

	assert.Equal(t, os.Getenv("HOMEDRIVE"), homeDrive, "output")
	assert.Equal(t, os.Getenv("HOMEPATH"), homePath, "output")
}
