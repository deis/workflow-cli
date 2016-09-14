package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) DomainsList(string, int) error {
	return nil
}

func (d FakeDeisCmd) DomainsAdd(string, string) error {
	return nil
}

func (d FakeDeisCmd) DomainsRemove(string, string) error {
	return nil
}

func TestDomains(t *testing.T) {
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
			args: []string{"domains:add", "example.com"},
		},
		{
			args: []string{"domains:list"},
		},
		{
			args: []string{"domains:remove", "example.com"},
		},
		{
			args: []string{"domains"},
		},
	}

	for _, c := range cases {
		err = Domains(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
