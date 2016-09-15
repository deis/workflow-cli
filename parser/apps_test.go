package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDeisCmd) AppCreate(string, string, string, bool) error {
	return errors.New("apps:create")
}

func (d FakeDeisCmd) AppsList(int) error {
	return errors.New("apps:list")
}

func (d FakeDeisCmd) AppInfo(string) error {
	return errors.New("apps:info")
}

func (d FakeDeisCmd) AppOpen(string) error {
	return errors.New("apps:open")
}

func (d FakeDeisCmd) AppLogs(string, int) error {
	return errors.New("apps:logs")
}

func (d FakeDeisCmd) AppRun(string, string) error {
	return errors.New("apps:run")
}

func (d FakeDeisCmd) AppDestroy(string, string) error {
	return errors.New("apps:destroy")
}

func (d FakeDeisCmd) AppTransfer(string, string) error {
	return errors.New("apps:transfer")
}

func TestApps(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDeisCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"apps:create"},
			expected: "",
		},
		{
			args:     []string{"apps:list"},
			expected: "",
		},
		{
			args:     []string{"apps:info"},
			expected: "",
		},
		{
			args:     []string{"apps:open"},
			expected: "",
		},
		{
			args:     []string{"apps:logs"},
			expected: "",
		},
		{
			args:     []string{"apps:logs", "--lines=1"},
			expected: "",
		},
		{
			args:     []string{"apps:run", "ls"},
			expected: "",
		},
		{
			args:     []string{"apps:destroy"},
			expected: "",
		},
		{
			args:     []string{"apps:transfer", "test-user"},
			expected: "",
		},
		{
			args:     []string{"apps"},
			expected: "apps:list",
		},
	}

	// For each case, check that calling the route with the arguments
	// returns the expected error, which is args[0] if not provided.
	for _, c := range cases {
		var expected string
		if c.expected == "" {
			expected = c.args[0]
		} else {
			expected = c.expected
		}
		err = Apps(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
