package parser

import (
	"bytes"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) CertsList(int, time.Time) error {
	return nil
}

func (d FakeDeisCmd) CertAdd(string, string, string) error {
	return nil
}

func (d FakeDeisCmd) CertRemove(string) error {
	return nil
}

func (d FakeDeisCmd) CertInfo(string) error {
	return nil
}

func (d FakeDeisCmd) CertAttach(string, string) error {
	return nil
}

func (d FakeDeisCmd) CertDetach(string, string) error {
	return nil
}

func TestCerts(t *testing.T) {
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
			args: []string{"certs:list"},
		},
		{
			args: []string{"certs:add", "name", "cert", "key"},
		},
		{
			args: []string{"certs:remove", "name"},
		},
		{
			args: []string{"certs:info", "name"},
		},
		{
			args: []string{"certs:attach", "name", "example.com"},
		},
		{
			args: []string{"certs:detach", "name", "example.com"},
		},
	}

	for _, c := range cases {
		err = Certs(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
