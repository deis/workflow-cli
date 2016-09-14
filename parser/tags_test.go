package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) TagsList(string) error {
	return nil
}

func (d FakeDeisCmd) TagsSet(string, []string) error {
	return nil
}

func (d FakeDeisCmd) TagsUnset(string, []string) error {
	return nil
}

func TestTags(t *testing.T) {
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
			args: []string{"tags:list"},
		},
		{
			args: []string{"tags:set", "environ", "prod"},
		},
		{
			args: []string{"tags:unset", "environ"},
		},
		{
			args: []string{"tags"},
		},
	}

	for _, c := range cases {
		err = Tags(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
