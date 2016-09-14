package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) GitRemote(string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) GitRemove(string) error {
	return nil
}

func TestGit(t *testing.T) {
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
			args: []string{"git:remote"},
		},
		{
			args: []string{"git:remove"},
		},
		{
			args: []string{"git"},
		},
	}

	for _, c := range cases {
		err = Git(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
