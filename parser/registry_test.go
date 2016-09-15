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

func (d FakeDeisCmd) RegistryList(string) error {
	return errors.New("registry:list")
}

func (d FakeDeisCmd) RegistrySet(string, []string) error {
	return errors.New("registry:set")
}

func (d FakeDeisCmd) RegistryUnset(string, []string) error {
	return errors.New("registry:unset")
}

func TestRegistry(t *testing.T) {
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
			args:     []string{"registry:list"},
			expected: "",
		},
		{
			args:     []string{"registry:set", "username", "value"},
			expected: "",
		},
		{
			args:     []string{"registry:unset", "username"},
			expected: "",
		},
		{
			args:     []string{"registry"},
			expected: "registry:list",
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
		err = Registry(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
