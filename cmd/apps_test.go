package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/arschles/assert"

	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/pkg/testutil"
)

type expandURLCases struct {
	Input    string
	Expected string
}

func TestExpandUrl(t *testing.T) {
	checks := []expandURLCases{
		{
			Input:    "test.com",
			Expected: "test.com",
		},
		{
			Input:    "test",
			Expected: "test.foo.com",
		},
	}

	for _, check := range checks {
		out := expandURL("deis.foo.com", check.Input)

		if out != check.Expected {
			t.Errorf("Expected %s, Got %s", check.Expected, out)
		}
	}
}

func TestRemoteExists(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "owner": "jkirk",
    "id": "foo",
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`)
	})

	// create a remote first before running apps:create
	dir, err := ioutil.TempDir("", "apps")
	assert.NoErr(t, err)

	defer os.RemoveAll(dir)

	assert.NoErr(t, os.Chdir(dir))

	assert.NoErr(t, git.Init(git.DefaultCmd))
	assert.NoErr(t, git.CreateRemote(git.DefaultCmd, "localhost", "deis", "appname"))

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.AppCreate("foo", "", "deis", false)

	assert.Equal(t, err.Error(), `A git remote with the name deis already exists. To overwrite this remote run:
deis git:remote --force --remote deis --app foo`,
		"output")
}
