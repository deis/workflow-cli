package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) AutoscaleList(string) error {
	return nil
}

func (d FakeDeisCmd) AutoscaleSet(string, string, int, int, int) error {
	return nil
}

func (d FakeDeisCmd) AutoscaleUnset(string, string) error {
	return nil
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

	cases := []struct {
		args []string
	}{
		{
			args: []string{"autoscale:list"},
		},
		{
			args: []string{"autoscale:set", "web/cmd", "--min=1", "--max=3", "--cpu-percent=50"},
		},
		{
			args: []string{"autoscale:unset", "web/cmd"},
		},
		{
			args: []string{"autoscale"},
		},
	}

	for _, c := range cases {
		err = Autoscale(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
