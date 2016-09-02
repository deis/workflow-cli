package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/pkg/testutil"
	"github.com/deis/workflow-cli/version"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("DEIS_API_VERSION", "1234")
		w.WriteHeader(200)
	})

	err = cmdr.Version(true)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), fmt.Sprintf(`Workflow CLI Version:            %s
Workflow CLI API Version:        %s
Workflow Controller API Version: 1234
`, version.Version, deis.APIVersion), "output")

	b.Reset()
	err = cmdr.Version(false)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), version.Version+"\n", "output")
}
