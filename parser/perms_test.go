package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) PermsList(string, bool, int) error {
	return nil
}

func (d FakeDeisCmd) PermCreate(string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) PermDelete(string, string, bool) error {
	return nil
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

	cases := []struct {
		args []string
	}{
		{
			args: []string{"perms:list"},
		},
		{
			args: []string{"perms:create", "test-user"},
		},
		{
			args: []string{"perms:delete", "test-user"},
		},
	}

	for _, c := range cases {
		err = Perms(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
