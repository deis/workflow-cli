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

type parseLimitCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseLimit(t *testing.T) {
	t.Parallel()

	cases := []parseLimitCase{
		{"web=2G", "web", "2G", false, ""},
		{"=1", "", "", true, `=1 doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`},
		{"web=", "", "", true, `web= doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`},
		{"1=", "", "", true, `1= doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`},
		{"web=G", "", "", true, `web=G doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`},
	}

	for _, check := range cases {
		key, value, err := parseLimit(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoErr(t, err)
			assert.Equal(t, key, check.Key, "key")
			assert.Equal(t, value, check.Value, "value")
		}
	}
}

type parseLimitsCase struct {
	Input         []string
	ExpectedMap   map[string]interface{}
	ExpectedError bool
	ExpectedMsg   string
}

func TestLimitTags(t *testing.T) {
	t.Parallel()

	cases := []parseLimitsCase{
		{[]string{"web=1G", "worker=2"}, map[string]interface{}{"web": "1G", "worker": "2"}, false, ""},
		{[]string{"foo=", "web=1G"}, nil, true, `foo= doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`},
	}

	for _, check := range cases {
		actual, err := parseLimits(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoErr(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "map")
		}
	}
}

func TestLimitsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/enterprise/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "enterprise",
			"values": {},
			"memory": {
				"web": "2G"
			},
			"cpu": {
				"web": "2",
				"worker": "1"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.LimitsList("enterprise")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== enterprise Limits

--- Memory
web     2G

--- CPU
web        2
worker     1
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
			"memory": {},
			"cpu": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
			}`)
	})
	b.Reset()

	err = cmdr.LimitsList("franklin")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== franklin Limits

--- Memory
Unlimited

--- CPU
Unlimited
`, "output")
}

func TestLimitsSet(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				CPU: map[string]interface{}{
					"web": "100m",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {},
			"cpu": {
				"web": "100m"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.LimitsSet("foo", []string{"web=100m"}, "cpu")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== foo Limits

--- Memory
Unlimited

--- CPU
web     100m
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Memory: map[string]interface{}{
					"web": "1G",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
			"memory": {
				"web": "1G"
			},
			"cpu": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsSet("franklin", []string{"web=1G"}, "memory")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== franklin Limits

--- Memory
web     1G

--- CPU
Unlimited
`, "output")
}

func TestLimitsUnset(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				Memory: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {},
			"cpu": {
				"web": "100m"
			},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.LimitsUnset("foo", []string{"web"}, "memory")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== foo Limits

--- Memory
Unlimited

--- CPU
web     100m
`, "output")

	server.Mux.HandleFunc("/v2/apps/franklin/config/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		if r.Method == "POST" {
			testutil.AssertBody(t, api.Config{
				CPU: map[string]interface{}{
					"web": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "bedison",
			"app": "franklin",
			"values": {},
			"memory": {
				"web": "1G"
			},
			"cpu": {},
			"tags": {},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})
	b.Reset()

	err = cmdr.LimitsUnset("franklin", []string{"web"}, "cpu")
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying limits... done

=== franklin Limits

--- Memory
web     1G

--- CPU
Unlimited
`, "output")
}
