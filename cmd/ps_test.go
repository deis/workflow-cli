package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/pkg/time"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestParseType(t *testing.T) {
	t.Parallel()

	// test RC pod name
	appID := "earthy-underdog"
	rcPod := "earthy-underdog-v2-cmd-8yngj"
	psType, psName := parseType(rcPod, appID)
	if psType != "cmd" || psName != rcPod {
		t.Errorf("type was not cmd (got %s) or psName was not %s (got %s)", psType, rcPod, psName)
	}

	// test Deployment pod name - they are longer due to hash
	appID = "nonfat-yearbook"
	deployPod := "nonfat-yearbook-cmd-2180299075-7na91"
	psType, psName = parseType(deployPod, appID)
	if psType != "cmd" || psName != deployPod {
		t.Errorf("type was not cmd (got %s) or psName was not %s (got %s)", psType, deployPod, psName)
	}

	// test type by itself
	psType, psName = parseType("cmd", "fake")
	if psType != "cmd" || psName != "" {
		t.Error("type was not cmd")
	}
}

func TestPrintProcesses(t *testing.T) {
	var b bytes.Buffer

	pods := []api.Pods{
		{
			Release: "v3",
			Name:    "benign-quilting-web-4084101150-c871y",
			Type:    "web",
			State:   "up",
			Started: time.Time{},
		},
		{
			Release: "v3",
			Name:    "benign-quilting-worker-4084101150-c871y",
			Type:    "worker",
			State:   "up",
			Started: time.Time{},
		},
	}

	printProcesses("appname", pods, &b)

	assert.Equal(t, b.String(), `=== appname Processes
--- web:
benign-quilting-web-4084101150-c871y up (v3)
--- worker:
benign-quilting-worker-4084101150-c871y up (v3)
`, "output")
}

func TestPsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"release": "v2",
					"type": "web",
					"name": "foo-web-4084101150-c871y",
					"state": "up",
					"started": "2016-02-13T00:47:52"
				}
			]
		}`)
	})

	err = cmdr.PsList("foo", -1)
	assert.NoErr(t, err)

	assert.Equal(t, b.String(), `=== foo Processes
--- web:
foo-web-4084101150-c871y up (v2)
`, "output")
}

type psTargetCases struct {
	Targets       []string
	ExpectedError bool
	ExpectedMap   map[string]int
	ExpectedMsg   string
}

func TestParsePsTargets(t *testing.T) {
	t.Parallel()

	cases := []psTargetCases{
		{[]string{"test"}, true, nil, "'test' does not match the pattern 'type=num', ex: web=2\n"},
		{[]string{"test=a"}, true, nil, "'test=a' does not match the pattern 'type=num', ex: web=2\n"},
		{[]string{"test="}, true, nil, "'test=' does not match the pattern 'type=num', ex: web=2\n"},
		{[]string{"test=2"}, false, map[string]int{"test": 2}, ""},
	}

	for _, check := range cases {
		actual, err := parsePsTargets(check.Targets)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoErr(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "error")
		}
	}
}

func TestPsScale(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}
	err = cmdr.PsScale("foo", []string{"test"})
	assert.Equal(t, err.Error(), "'test' does not match the pattern 'type=num', ex: web=2\n", "error")

	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"release": "v2",
					"type": "web",
					"name": "foo-web-4084101150-c871y",
					"state": "up",
					"started": "2016-02-13T00:47:52"
				}
			]
		}`)
	})

	server.Mux.HandleFunc("/v2/apps/foo/scale/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]int{"web": 1}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()
	err = cmdr.PsScale("foo", []string{"web=1"})
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Scaling processes... but first, coffee!
done in 0s
=== foo Processes
--- web:
foo-web-4084101150-c871y up (v2)
`, "output")
}

func TestPsRestart(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `[
		{
				"release": "v2",
				"type": "web",
				"name": "foo-web-4084101150-c871y",
				"state": "up",
				"started": "2016-02-13T00:47:52"
		}
]`)
	})

	b.Reset()
	err = cmdr.PsRestart("foo", "")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Restarting processes... but first, coffee!
done in 0s
=== foo Processes
--- web:
foo-web-4084101150-c871y up (v2)
`, "output")

	server.Mux.HandleFunc("/v2/apps/coolapp/pods/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `[]`)
	})

	b.Reset()

	err = cmdr.PsRestart("coolapp", "")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Restarting processes... but first, coffee!
Could not find any processes to restart
`, "output")

	server.Mux.HandleFunc("/v2/apps/testapp/pods/web/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `[
			{
				"release": "v2",
				"type": "web",
				"name": "testapp-web-4084101150-c871y",
				"state": "up",
				"started": "2016-02-13T00:47:52"
			}
		]`)
	})

	b.Reset()

	err = cmdr.PsRestart("testapp", "web")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Restarting processes... but first, coffee!
done in 0s
=== testapp Processes
--- web:
testapp-web-4084101150-c871y up (v2)
`, "output")

	server.Mux.HandleFunc("/v2/apps/newapp/pods/web/newapp-web-4084101150-c871y/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `[
			{
				"release": "v2",
				"type": "web",
				"name": "newapp-web-4084101150-c871y",
				"state": "up",
				"started": "2016-02-13T00:47:52"
			}
		]`)
	})

	b.Reset()

	err = cmdr.PsRestart("newapp", "newapp-web-4084101150-c871y")
	assert.NoErr(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Restarting processes... but first, coffee!
done in 0s
=== newapp Processes
--- web:
newapp-web-4084101150-c871y up (v2)
`, "output")

	server.Mux.HandleFunc("/v2/apps/newapp/pods/ghost/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{
			"detail": "Container type ghost does not exist in application"
		}`)
	})

	b.Reset()

	err = cmdr.PsRestart("newapp", "ghost")
	assert.Equal(t, err.Error(), "Could not find process type ghost in app newapp", "error")
}
