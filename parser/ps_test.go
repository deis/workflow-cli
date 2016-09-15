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

func (d FakeDeisCmd) PsList(string, int) error {
	return errors.New("ps:list")
}

func (d FakeDeisCmd) PsScale(string, []string) error {
	return errors.New("ps:scale")
}

func (d FakeDeisCmd) PsRestart(string, string) error {
	return errors.New("ps:restart")
}

func TestPs(t *testing.T) {
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
			args:     []string{"ps:list"},
			expected: "",
		},
		{
			args:     []string{"ps:restart", "web"},
			expected: "",
		},
		{
			args:     []string{"ps:scale", "web", "5"},
			expected: "",
		},
		{
			args:     []string{"ps:list"},
			expected: "",
		},
		{
			args:     []string{"ps"},
			expected: "ps:list",
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
		err = Ps(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
