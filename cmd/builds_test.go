package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestParseProcfile(t *testing.T) {
	t.Parallel()

	procMap, err := parseProcfile([]byte(`web: ./test
foo: test --test
`))
	assert.NoErr(t, err)
	assert.Equal(t, procMap, map[string]string{"web": "./test", "foo": "test --test"}, "map")

	_, err = parseProcfile([]byte(`web: ./test
foo
`))

	assert.ExistsErr(t, err, "yaml")
}

func TestBuildsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"app": "",
					"created": "2014-01-01T00:00:00UTC",
					"dockerfile": "",
					"image": "",
					"owner": "",
					"procfile": {},
					"sha": "",
					"updated": "",
					"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
				},
				{
					"app": "",
					"created": "2014-01-05T00:00:00UTC",
					"dockerfile": "",
					"image": "",
					"owner": "",
					"procfile": {},
					"sha": "",
					"updated": "",
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3"
				}
			]
		}`)
	})

	err = cmdr.BuildsList("foo", -1)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== foo Builds
de1bf5b5-4a72-4f94-a10c-d2a3741cdf75 2014-01-01T00:00:00UTC
c4aed81c-d1ca-4ff1-ab89-d2151264e1a3 2014-01-05T00:00:00UTC
`, "output")
}

func TestBuildsCreate(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	// Create a new temporary directory and change to it.
	name, err := ioutil.TempDir("", "client")
	assert.NoErr(t, err)
	err = os.Chdir(name)
	assert.NoErr(t, err)

	server.Mux.HandleFunc("/v2/apps/enterprise/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "ncc/1701:A",
		}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	err = cmdr.BuildsCreate("enterprise", "ncc/1701:A", "")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/bradbury/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "nx/72307:latest",
			Procfile: map[string]string{
				"web":  "./drive",
				"warp": "./warp 8",
			},
		}, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})
	b.Reset()

	err = cmdr.BuildsCreate("bradbury", "nx/72307:latest", `web: ./drive
warp: ./warp 8
`)
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

	server.Mux.HandleFunc("/v2/apps/franklin/builds/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.CreateBuildRequest{
			Image: "nx/326:latest",
			Procfile: map[string]string{
				"web":  "./drive",
				"warp": "./warp 8",
			},
		}, r)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})
	b.Reset()

	err = ioutil.WriteFile("Procfile", []byte(`web: ./drive
warp: ./warp 8
`), os.ModePerm)
	assert.NoErr(t, err)

	err = cmdr.BuildsCreate("franklin", "nx/326:latest", "")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Creating build... done\n", "output")

}
