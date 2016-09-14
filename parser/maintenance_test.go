package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) MaintenanceInfo(string) error {
	return nil
}

func (d FakeDeisCmd) MaintenanceEnable(string) error {
	return nil
}

func (d FakeDeisCmd) MaintenanceDisable(string) error {
	return nil
}

func TestMaintenance(t *testing.T) {
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
			args: []string{"maintenance:info"},
		},
		{
			args: []string{"maintenance:on"},
		},
		{
			args: []string{"maintenance:off"},
		},
		{
			args: []string{"maintenance"},
		},
	}

	for _, c := range cases {
		err = Maintenance(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
