package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) WhitelistAdd(string, string) error {
	return nil
}

func (d FakeDeisCmd) WhitelistList(string) error {
	return nil
}

func (d FakeDeisCmd) WhitelistRemove(string, string) error {
	return nil
}

func TestWhitelist(t *testing.T) {
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
			args: []string{"whitelist:add", "1.2.3.4"},
		},
		{
			args: []string{"whitelist:list"},
		},
		{
			args: []string{"whitelist:remove", "1.2.3.4"},
		},
	}

	for _, c := range cases {
		err = Whitelist(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
