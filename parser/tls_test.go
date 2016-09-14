package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) TLSInfo(string) error {
	return nil
}

func (d FakeDeisCmd) TLSEnable(string) error {
	return nil
}

func (d FakeDeisCmd) TLSDisable(string) error {
	return nil
}

func TestTLS(t *testing.T) {
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
			args: []string{"tls:info"},
		},
		{
			args: []string{"tls:enable"},
		},
		{
			args: []string{"tls:disable"},
		},
		{
			args: []string{"tls"},
		},
	}

	for _, c := range cases {
		err = TLS(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
