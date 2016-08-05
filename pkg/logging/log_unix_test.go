// +build linux darwin

package logging

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
)

type colorsTestCase struct {
	Input    string
	Expected string
}

func TestChooseColor(t *testing.T) {
	t.Parallel()

	colors := []colorsTestCase{
		{"INFO", "\033[35m"},
		{"a", "\033[32m"},
		{"b", "\033[33m"},
		{"c", "\033[34m"},
		{"d", "\033[39m"},
		{"e", "\033[36m"},
		{"f", "\033[31m"},
	}
	for _, color := range colors {
		assert.Equal(t, chooseColor(color.Input), color.Expected, "color")
	}
}

func TestPrintLogs(t *testing.T) {
	var b bytes.Buffer
	PrintLog(&b, "INFO [test]: testing")
	assert.Equal(t, b.String(), "\033[35mINFO [test]: testing\033[0m\n", "log line")
	b.Reset()
	// Regression test for https://github.com/deis/deis/issues/4420
	PrintLog(&b, "\nDone preparing production files\n\n\u001b[4mRunning \"concat:plugins\" (concat) task\u001b[24m\n")
	assert.Equal(t, b.String(),
		"\033[31m\nDone preparing production files\n\n\u001b[4mRunning \"concat:plugins\" (concat) task\u001b[24m\n\033[0m\n", "log line")
}
