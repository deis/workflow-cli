package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestRoutingInfo(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/rivendell/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "elrond",
			"app": "rivendell",
			"maintenance": true,
			"routable": true,
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	err = cmdr.RoutingInfo("rivendell")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Routing is enabled.\n", "output")

	server.Mux.HandleFunc("/v2/apps/mordor/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "sauron",
			"app": "mordor",
			"maintenance": true,
			"routable": false,
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.RoutingInfo("mordor")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Routing is disabled.\n", "output")

	// test that no routable field doesn't trigger a panic
	server.Mux.HandleFunc("/v2/apps/gondor/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "aragorn",
			"app": "gondor",
			"maintenance": true,
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.RoutingInfo("gondor")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Routing is enabled.\n", "output")
}

func TestRoutingEnable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/lothlorien/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.AppSettings{Routable: api.NewRoutable()}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.RoutingEnable("lothlorien")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Enabling routing for lothlorien... done\n", "output")
}

func TestRoutingDisable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/bree/settings/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		routable := false
		testutil.AssertBody(t, api.AppSettings{Routable: &routable}, r)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.RoutingDisable("bree")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Disabling routing for bree... done\n", "output")
}
