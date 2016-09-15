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

func (d FakeDeisCmd) PermsList(string, bool, int) error {
	return errors.New("perms:list")
}

func (d FakeDeisCmd) PermCreate(string, string, bool) error {
	return errors.New("perms:create")
}

func (d FakeDeisCmd) PermDelete(string, string, bool) error {
	return errors.New("perms:delete")
}

func TestPerms(t *testing.T) {
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
			args:     []string{"perms:list"},
			expected: "",
		},
		{
			args:     []string{"perms:create", "test-user"},
			expected: "",
		},
		{
			args:     []string{"perms:delete", "test-user"},
			expected: "",
		},
		{
			args:     []string{"perms"},
			expected: "perms:list",
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
		err = Perms(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
