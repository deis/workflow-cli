// +build windows

package settings

import "os"

// FindHome returns the HOME directory of the current user
func FindHome() string {
	return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
}

// SetHome sets the HOME directory of the current user
func SetHome(path string) {
	os.Setenv("HOMEDRIVE", path[:2])
	os.Setenv("HOMEPATH", path[2:])
}
