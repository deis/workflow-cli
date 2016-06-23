// +build windows

package cmd

// Windows doesn't support colored terminal output, so this is just a dummy method.
func chooseColor(input string) string {
	return input
}
