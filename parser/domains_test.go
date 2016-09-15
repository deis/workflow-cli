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

func (d FakeDeisCmd) DomainsList(string, int) error {
	return errors.New("domains:list")
}

func (d FakeDeisCmd) DomainsAdd(string, string) error {
	return errors.New("domains:add")
}

func (d FakeDeisCmd) DomainsRemove(string, string) error {
	return errors.New("domains:remove")
}

func TestDomains(t *testing.T) {
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
			args:     []string{"domains:add", "example.com"},
			expected: "",
		},
		{
			args:     []string{"domains:list"},
			expected: "",
		},
		{
			args:     []string{"domains:remove", "example.com"},
			expected: "",
		},
		{
			args:     []string{"domains"},
			expected: "domains:list",
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
		err = Domains(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
