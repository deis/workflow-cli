package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) PsList(string, int) error {
	return nil
}

func (d FakeDeisCmd) PsScale(string, []string) error {
	return nil
}

func (d FakeDeisCmd) PsRestart(string, string) error {
	return nil
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

	cases := []struct {
		args []string
	}{
		{
			args: []string{"ps:list"},
		},
		{
			args: []string{"ps:restart", "web"},
		},
		{
			args: []string{"ps:scale", "web", "5"},
		},
	}

	for _, c := range cases {
		err = Ps(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
