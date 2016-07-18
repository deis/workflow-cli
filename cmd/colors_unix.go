// +build linux darwin

package cmd

import "fmt"

const colorStringEscape = "\033[3%dm"

// Choose an ANSI color by converting a string to an int.
//
// Color 5, magenta, reserved for controller log messages.
//
// Colors 0 and 7, black and white, are skipped because they are likely to be unreadable
// against terminal backgrounds. Instead, 9, the default text color is used, likely to be the color
// of the two that is readable against the terminal background.
func chooseColor(input string) string {
	if input == "INFO" {
		return fmt.Sprintf(colorStringEscape, 5)
	}

	var sum uint8

	for _, char := range []byte(input) {
		sum += uint8(char)
	}

	// Eight possible terminal colors, but Black and White are excluded
	color := (sum % 6) + 1

	// Color 5 is reserved, replace it with default text color.
	if color == 5 {
		color = 9
	}

	return fmt.Sprintf(colorStringEscape, color)
}
