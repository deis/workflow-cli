package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) BuildsList(string, int) error {
	return nil
}

func (d FakeDeisCmd) BuildsCreate(string, string, string) error {
	return nil
}

func TestBuilds(t *testing.T) {
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
			args: []string{"builds:list"},
		},
		{
			args: []string{"builds:create", "deis/example-go:latest"},
		},
		{
			args: []string{"builds"},
		},
	}

	for _, c := range cases {
		err = Builds(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
