package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) KeysList(int) error {
	return nil
}

func (d FakeDeisCmd) KeyRemove(string) error {
	return nil
}

func (d FakeDeisCmd) KeyAdd(string) error {
	return nil
}

func TestKeys(t *testing.T) {
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
			args: []string{"keys:list"},
		},
		{
			args: []string{"keys:add", "key"},
		},
		{
			args: []string{"keys:remove", "key"},
		},
	}

	for _, c := range cases {
		err = Keys(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
