package parser

import (
	"bytes"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func (d FakeDeisCmd) Register(string, string, string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) Login(string, string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) Logout() error {
	return nil
}

func (d FakeDeisCmd) Passwd(string, string, string) error {
	return nil
}

func (d FakeDeisCmd) Cancel(string, string, bool) error {
	return nil
}

func (d FakeDeisCmd) Whoami(bool) error {
	return nil
}

func (d FakeDeisCmd) Regenerate(string, bool) error {
	return nil
}

// UsersList lists users registered with the controller.
func (d FakeDeisCmd) UsersList(results int) error {
	return nil
}

func TestUsers(t *testing.T) {
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
			args: []string{"users:list"},
		},
		{
			args: []string{"users"},
		},
	}

	for _, c := range cases {
		err = Users(c.args, cmdr)
		assert.NoErr(t, err)
	}
}
