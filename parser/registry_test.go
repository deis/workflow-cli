package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) RegistryList(string) error {
	return nil
}

func (d FakeDeisCmd) RegistrySet(string, []string) error {
	return nil
}

func (d FakeDeisCmd) RegistryUnset(string, []string) error {
	return nil
}

func TestRegistry(t *testing.T) {
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
			args: []string{"registry:list"},
		},
		{
			args: []string{"registry:set", "username", "value"},
		},
		{
			args: []string{"registry:unset", "username"},
		},
		{
			args: []string{"registry"},
		},
	}

	for _, c := range cases {
		err = Registry(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
