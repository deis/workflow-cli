// +build linux darwin

package cmd

import "fmt"

// Choose an ANSI color by converting a string to an int.
func chooseColor(input string) string {
	var sum uint8

	for _, char := range []byte(input) {
		sum += uint8(char)
	}

	// Seven possible terminal colors
	color := (sum % 7) + 1

	if color == 7 {
		color = 9
	}

	return fmt.Sprintf("\033[3%dm", color)
}
