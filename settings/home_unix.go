// +build linux darwin

package settings

import (
	"os"
)

// FindHome returns the HOME directory of the current user
func FindHome() string {
	return os.Getenv("HOME")
}
