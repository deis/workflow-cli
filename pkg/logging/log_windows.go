// +build windows

package logging

import (
	"fmt"
	"io"
)

// PrintLog prints each log line. Windows doesn't support ansi escape codes,
// so color is not printed.
func PrintLog(out io.Writer, log string) {
	fmt.Fprintln(out, log)
}
