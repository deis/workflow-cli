package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) RoutingInfo(string) error {
	return nil
}

func (d FakeDeisCmd) RoutingEnable(string) error {
	return nil
}

func (d FakeDeisCmd) RoutingDisable(string) error {
	return nil
}

func TestRouting(t *testing.T) {
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
			args: []string{"routing:info"},
		},
		{
			args: []string{"routing:enable"},
		},
		{
			args: []string{"routing:disable"},
		},
	}

	for _, c := range cases {
		err = Routing(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
