package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) ConfigList(string, bool) error {
	return nil
}

func (d FakeDeisCmd) ConfigSet(string, []string) error {
	return nil
}

func (d FakeDeisCmd) ConfigUnset(string, []string) error {
	return nil
}

func (d FakeDeisCmd) ConfigPull(string, bool, bool) error {
	return nil
}

func (d FakeDeisCmd) ConfigPush(string, string) error {
	return nil
}

func TestConfig(t *testing.T) {
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
			args: []string{"config:list"},
		},
		{
			args: []string{"config:set", "var=value"},
		},
		{
			args: []string{"config:unset", "var"},
		},
		{
			args: []string{"config:pull"},
		},
		{
			args: []string{"config:push"},
		},
	}

	for _, c := range cases {
		err = Config(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
