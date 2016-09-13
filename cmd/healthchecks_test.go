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

func TestPrintHealthCheck(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	testHealthCheck := api.Healthchecks{}
	cmdr.printHealthCheck(testHealthCheck)
	assert.Equal(t, b.String(), "--- Liveness\nNo liveness probe configured.\n\n--- Readiness\nNo readiness probe configured.\n", "healthcheck")
	b.Reset()
	testHealthCheck["livenessProbe"] = &api.Healthcheck{}
	testHealthCheck["readinessProbe"] = &api.Healthcheck{}
	cmdr.printHealthCheck(testHealthCheck)
	assert.Equal(t, b.String(), "--- Liveness\nInitial Delay (seconds): 0\nTimeout (seconds): 0\nPeriod (seconds): 0\nSuccess Threshold: 0\nFailure Threshold: 0\nExec Probe: N/A\nHTTP GET Probe: N/A\nTCP Socket Probe: N/A\n\n--- Readiness\nInitial Delay (seconds): 0\nTimeout (seconds): 0\nPeriod (seconds): 0\nSuccess Threshold: 0\nFailure Threshold: 0\nExec Probe: N/A\nHTTP GET Probe: N/A\nTCP Socket Probe: N/A\n", "healthcheck")
}

func TestHealthchecksList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web/cmd": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksList("foo", "web/cmd")
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== foo Healthchecks

web/cmd:
--- Liveness
Initial Delay (seconds): 50
Timeout (seconds): 50
Period (seconds): 10
Success Threshold: 1
Failure Threshold: 3
Exec Probe: N/A
HTTP GET Probe: Path="/" Port=80 HTTPHeaders=[]
TCP Socket Probe: N/A

--- Readiness
No readiness probe configured.
`, "output")
}

func TestHealthchecksSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {
    "web/cmd": {
      "livenessProbe": {
        "initialDelaySeconds": 50,
        "timeoutSeconds": 50,
        "periodSeconds": 10,
        "failureThreshold": 3,
        "httpGet": {
          "port": 80,
          "path": "/"
        },
        "successThreshold": 1
      }
    }
  },
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksSet("foo", "liveness", "web/cmd", &api.Healthcheck{})
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Applying liveness healthcheck... done

=== foo Healthchecks

web/cmd:
--- Liveness
Initial Delay (seconds): 50
Timeout (seconds): 50
Period (seconds): 10
Success Threshold: 1
Failure Threshold: 3
Exec Probe: N/A
HTTP GET Probe: Path="/" Port=80 HTTPHeaders=[]
TCP Socket Probe: N/A

--- Readiness
No readiness probe configured.
`, "output")
}

func TestHealthchecksUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "uuid": "c039a380-6068-4511-b35a-535a73b86ef5",
  "app": "foo",
  "owner": "bar",
  "values": {},
  "memory": {},
  "cpu": {},
  "tags": {},
  "registry": {},
  "healthcheck": {},
  "created": "2016-09-12T22:20:14Z",
  "updated": "2016-09-12T22:20:14Z"
}`)
	})

	err = cmdr.HealthchecksUnset("foo", "web/cmd", []string{"web/cmd"})
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Removing healthchecks... done

=== foo Healthchecks

web/cmd:
--- Liveness
No liveness probe configured.

--- Readiness
No readiness probe configured.
`, "output")
}
