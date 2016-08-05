// +build linux darwin

package logging

import (
	"fmt"
	"io"
	"strings"

	"github.com/deis/pkg/prettyprint"
)

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

// PrintLog prints a log line with a color matched to its category.
func PrintLog(out io.Writer, log string) {
	category := "unknown"
	parts := strings.Split(strings.Split(log, " -- ")[0], " ")
	category = parts[0]
	colorVars := map[string]string{
		"Color": chooseColor(category),
		"Log":   log,
	}
	fmt.Fprintln(out, prettyprint.ColorizeVars("{{.V.Color}}{{.V.Log}}{{.C.Default}}", colorVars))
}
