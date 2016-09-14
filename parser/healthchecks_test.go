package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) HealthchecksList(string, string) error {
	return nil
}

func (d FakeDeisCmd) HealthchecksSet(string, string, string, *api.Healthcheck) error {
	return nil
}

func (d FakeDeisCmd) HealthchecksUnset(string, string, []string) error {
	return nil
}

func TestHealthchecks(t *testing.T) {
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
			args: []string{"healthchecks:list"},
		},
		{
			args: []string{"healthchecks:set", "liveness", "httpGet", "80"},
		},
		{
			args: []string{"healthchecks:unset", "liveness"},
		},
	}

	for _, c := range cases {
		err = Healthchecks(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
