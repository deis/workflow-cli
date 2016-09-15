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

func (d FakeDeisCmd) ConfigList(string, bool) error {
	return errors.New("config:list")
}

func (d FakeDeisCmd) ConfigSet(string, []string) error {
	return errors.New("config:set")
}

func (d FakeDeisCmd) ConfigUnset(string, []string) error {
	return errors.New("config:unset")
}

func (d FakeDeisCmd) ConfigPull(string, bool, bool) error {
	return errors.New("config:pull")
}

func (d FakeDeisCmd) ConfigPush(string, string) error {
	return errors.New("config:push")
}

func TestConfig(t *testing.T) {
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
			args:     []string{"config:list"},
			expected: "",
		},
		{
			args:     []string{"config:set", "var=value"},
			expected: "",
		},
		{
			args:     []string{"config:unset", "var"},
			expected: "",
		},
		{
			args:     []string{"config:pull"},
			expected: "",
		},
		{
			args:     []string{"config:push"},
			expected: "",
		},
		{
			args:     []string{"config"},
			expected: "config:list",
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
		err = Config(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
