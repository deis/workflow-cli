// +build windows

package settings

import (
	"os"
	"testing"

	"github.com/arschles/assert"
)

// TestFindHome ensures the correct home directory is returned by FindHome().
func TestFindHome(t *testing.T) {
	previousHomedrive := os.Getenv("HOMEDRIVE")
	previousHomepath := os.Getenv("HOMEPATH")

	homedrive := "C:"
	homepath := "/a/b/c"
	os.Setenv("HOMEDRIVE", homedrive)
	os.Setenv("HOMEPATH", homepath)
	assert.Equal(t, FindHome(), homedrive+homepath, "output")

	// Reset vars.
	os.Setenv("HOMEDRIVE", previousHomedrive)
	os.Setenv("HOMEPATH", previousHomepath)
}

// TestSetHome ensures the correct env vars are set when SetHome() is called.
func TestSetHome(t *testing.T) {
	previousHomedrive := os.Getenv("HOMEDRIVE")
	previousHomepath := os.Getenv("HOMEPATH")

	homeDrive := "D:"
	homePath := "/e/f/g"
	SetHome(homeDrive + homePath)

	assert.Equal(t, os.Getenv("HOMEDRIVE"), homeDrive, "output")
	assert.Equal(t, os.Getenv("HOMEPATH"), homePath, "output")

	// Reset vars.
	os.Setenv("HOMEDRIVE", previousHomedrive)
	os.Setenv("HOMEPATH", previousHomepath)
}
