package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestAuth(t *testing.T) {
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
			args: []string{"auth:register", server.Server.URL},
		},
		{
			args: []string{"auth:login", server.Server.URL},
		},
		{
			args: []string{"auth:login", server.Server.URL, "--ssl-verify=true"},
		},
		{
			args: []string{"auth:logout"},
		},
		{
			args: []string{"auth:passwd"},
		},
		{
			args: []string{"auth:whoami"},
		},
		{
			args: []string{"auth:cancel"},
		},
		{
			args: []string{"auth:regenerate"},
		},
	}

	for _, c := range cases {
		err = Auth(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
