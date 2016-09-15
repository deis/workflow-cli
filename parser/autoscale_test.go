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

func (d FakeDeisCmd) AutoscaleList(string) error {
	return errors.New("autoscale:list")
}

func (d FakeDeisCmd) AutoscaleSet(string, string, int, int, int) error {
	return errors.New("autoscale:set")
}

func (d FakeDeisCmd) AutoscaleUnset(string, string) error {
	return errors.New("autoscale:unset")
}

func TestAutoscale(t *testing.T) {
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
			args:     []string{"autoscale:list"},
			expected: "",
		},
		{
			args:     []string{"autoscale:set", "web/cmd", "--min=1", "--max=3", "--cpu-percent=50"},
			expected: "",
		},
		{
			args:     []string{"autoscale:unset", "web/cmd"},
			expected: "",
		},
		{
			args:     []string{"autoscale"},
			expected: "autoscale:list",
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
		err = Autoscale(c.args, cmdr)
		assert.Err(t, errors.New(expected), err)
	}
}
