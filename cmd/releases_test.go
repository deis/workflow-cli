package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestReleasesList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/releases/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
			{
					"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
					"app": "numenor",
					"owner": "nazgul",
					"created": "2016-08-22T17:40:16Z",
					"updated": "2016-08-22T17:40:16Z",
					"version": 2,
					"summary": "khamul added ANGMAR",
					"config": "3bb816b1-4fde-4b06-8afe-acd12f58a266",
					"build": null
				},
				{
					"app": "numenor",
					"build": null,
					"config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
					"created": "2014-01-01T00:00:00UTC",
					"owner": "nazgul",
					"summary": "nazgul created initial release",
					"updated": "2014-01-01T00:00:00UTC",
					"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
					"version": 1
				}
			]
		}`)
	})

	err = cmdr.ReleasesList("numenor", -1)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== numenor Releases
v2	2016-08-22T17:40:16Z	khamul added ANGMAR
v1	2014-01-01T00:00:00UTC	nazgul created initial release
`, "output")
}

func TestReleasesInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/releases/v2/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
			"app": "numenor",
			"owner": "nazgul",
			"created": "2016-08-22T17:40:16Z",
			"updated": "2016-08-22T17:40:16Z",
			"version": 2,
			"summary": "khamul added ANGMAR",
			"config": "3bb816b1-4fde-4b06-8afe-acd12f58a266",
			"build": null
		}`)
	})

	err = cmdr.ReleasesInfo("numenor", 2)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== numenor Release v2
config:   3bb816b1-4fde-4b06-8afe-acd12f58a266
owner:    nazgul
created:  2016-08-22T17:40:16Z
summary:  khamul added ANGMAR
updated:  2016-08-22T17:40:16Z
uuid:     c4aed81c-d1ca-4ff1-ab89-d2151264e1a3
`, "output")
}

func TestReleasesRollback(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/numenor/releases/rollback/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, body, []byte{}, "body")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"version": 5}`)
	})

	err = cmdr.ReleasesRollback("numenor", -1)
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Rolling back one release... done, v5\n", "output")

	server.Mux.HandleFunc("/v2/apps/angmar/releases/rollback/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.ReleaseRollback{Version: 3}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"version": 3}`)
	})

	b.Reset()

	err = cmdr.ReleasesRollback("angmar", 3)
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Rolling back to v3... done, v3\n", "output")
}
