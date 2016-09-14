package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) AppCreate(string, string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) AppsList(int) error {
	return nil
}

func (d FakeDeisCmd) AppInfo(string) error {
	return nil
}

func (d FakeDeisCmd) AppOpen(string) error {
	return nil
}

func (d FakeDeisCmd) AppLogs(string, int) error {
	return nil
}

func (d FakeDeisCmd) AppRun(string, string) error {
	return nil
}

func (d FakeDeisCmd) AppDestroy(string, string) error {
	return nil
}

func (d FakeDeisCmd) AppTransfer(string, string) error {
	return nil
}

func TestApps(t *testing.T) {
	t.Parallel()

	cases := []struct {
		args []string
	}{
		{
			args: []string{"apps:create"},
		},
		{
			args: []string{"apps:list"},
		},
		{
			args: []string{"apps:info"},
		},
		{
			args: []string{"apps:open"},
		},
		{
			args: []string{"apps:logs"},
		},
		{
			args: []string{"apps:logs", "--lines=1"},
		},
		{
			args: []string{"apps:run", "ls"},
		},
		{
			args: []string{"apps:destroy"},
		},
		{
			args: []string{"apps:transfer", "test-user"},
		},
		{
			args: []string{"apps"},
		},
	}

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDeisCmd{WOut: &b, ConfigFile: cf}

	for _, c := range cases {
		err = Apps(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
