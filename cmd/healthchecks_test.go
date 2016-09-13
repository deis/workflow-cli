package cmd

import (
	"bytes"
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
