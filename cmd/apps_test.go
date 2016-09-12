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
	"github.com/deis/workflow-cli/settings"
)

type expandURLCases struct {
	Input    string
	Expected string
}

func TestAppsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
			    {
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "lorem-ipsum",
					"owner": "dolar-sit-amet",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				},
				{
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "consectetur",
					"owner": "adipiscing",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				}
			]
		}`)
	})

	err = cmdr.AppsList(-1)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== Apps
lorem-ipsum
consectetur
`, "output")
}

func TestAppsListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
			    {
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"id": "lorem-ipsum",
					"owner": "dolar-sit-amet",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"structure": {
						"cmd": 1
					}
				}
			]
		}`)
	})

	err = cmdr.AppsList(1)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== Apps (1 of 2)
lorem-ipsum
`, "output")
}

func TestAppsInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/pods/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "results": [
      {
        "state": "up",
        "started": "2016-08-22T17:42:16Z",
        "name": "lorem-ipsum-cmd-1911796442-48b58",
        "release": "v2",
        "type": "cmd"
      }
    ]
}`)
	})

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/domains/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
      {
         "owner": "dolar-sit-amet",
         "created": "2016-08-22T17:40:16Z",
         "updated": "2016-08-22T17:40:16Z",
         "app": "lorem-ipsum",
         "domain": "lorem-ipsum"
      }
    ]
}`)
	})

	s, err := settings.Load(cmdr.ConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	url, err := cmdr.appURL(s, "lorem-ipsum")
	if err != nil {
		t.Fatal(err)
	}

	if url == "" {
		url = fmt.Sprintf(noDomainAssignedMsg, "lorem-ipsum")
	}

	err = cmdr.AppInfo("lorem-ipsum")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== lorem-ipsum Application
updated:  2016-08-22T17:40:16Z
uuid:     c4aed81c-d1ca-4ff1-ab89-d2151264e1a3
created:  2016-08-22T17:40:16Z
url:      `+url+`
owner:    dolar-sit-amet
id:       lorem-ipsum

=== lorem-ipsum Processes
--- cmd:
lorem-ipsum-cmd-1911796442-48b58 up (v2)

=== lorem-ipsum Domains
lorem-ipsum

`, "output")
}

func TestAppDestroy(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})
	err = cmdr.AppDestroy("lorem-ipsum", "bad-confirm-string")
	assert.Equal(t, err.Error(), `App lorem-ipsum does not match confirm bad-confirm-string, aborting.`, "output")

	err = cmdr.AppDestroy("lorem-ipsum", "lorem-ipsum")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Destroying lorem-ipsum...
done in 0s
`, "output")
}

func TestAppTransfer(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lorem-ipsum/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "id": "lorem-ipsum",
    "owner": "dolar-sit-amet",
    "structure": {
      "cmd": 1
    },
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z"
}`)
	})

	err = cmdr.AppTransfer("lorem-ipsum", "test-user")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Transferring lorem-ipsum to test-user... done
`, "output")
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
