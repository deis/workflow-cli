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

func (d FakeDeisCmd) LimitsList(string) error {
	return errors.New("limits:list")
}

func (d FakeDeisCmd) LimitsSet(string, []string, string) error {
	return errors.New("limits:set")
}

func (d FakeDeisCmd) LimitsUnset(string, []string, string) error {
	return errors.New("limits:unset")
}

func TestLimits(t *testing.T) {
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
			args:     []string{"limits:list"},
			expected: "",
		},
		{
			args:     []string{"limits:set", "web=1G"},
			expected: "",
		},
		{
			args:     []string{"limits:set", "--cpu", "web=1"},
			expected: "",
		},
		{
			args:     []string{"limits:unset", "web"},
			expected: "",
		},
		{
			args:     []string{"limits:unset", "--cpu", "web"},
			expected: "",
		},
		{
			args:     []string{"limits"},
			expected: "limits:list",
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
		err = Limits(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
