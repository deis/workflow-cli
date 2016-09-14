package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) LimitsList(string) error {
	return nil
}

func (d FakeDeisCmd) LimitsSet(string, []string, string) error {
	return nil
}

func (d FakeDeisCmd) LimitsUnset(string, []string, string) error {
	return nil
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

	cases := []struct {
		args []string
	}{
		{
			args: []string{"limits:list"},
		},
		{
			args: []string{"limits:set", "web=1G"},
		},
		{
			args: []string{"limits:unset", "web"},
		},
	}

	for _, c := range cases {
		err = Limits(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
