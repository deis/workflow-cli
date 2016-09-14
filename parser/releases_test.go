package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) ReleasesList(string, int) error {
	return nil
}

func (d FakeDeisCmd) ReleasesInfo(string, int) error {
	return nil
}

func (d FakeDeisCmd) ReleasesRollback(string, int) error {
	return nil
}

func TestReleases(t *testing.T) {
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
			args: []string{"releases:list"},
		},
		{
			args: []string{"releases:info", "v1"},
		},
		{
			args: []string{"releases:rollback", "v1"},
		},
	}

	for _, c := range cases {
		err = Releases(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
