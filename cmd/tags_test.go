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

type parseTagCase struct {
	Input         string
	Key           string
	Value         string
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseTag(t *testing.T) {
	t.Parallel()

	cases := []parseTagCase{
		{"foo=bar", "foo", "bar", false, ""},
		{"test=1", "test", "1", false, ""},
		{"test", "", "", true, `test is invalid, Must be in format key=value
Examples: rack=1 evironment=production`},
		{"test=1=2", "", "", true, `test=1=2 is invalid, Must be in format key=value
Examples: rack=1 evironment=production`},
		{"test=", "", "", true, `test= is invalid, Must be in format key=value
Examples: rack=1 evironment=production`},
		{"=test", "", "", true, `=test is invalid, Must be in format key=value
Examples: rack=1 evironment=production`},
	}

	for _, check := range cases {
		key, value, err := parseTag(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoErr(t, err)
			assert.Equal(t, key, check.Key, "key")
			assert.Equal(t, value, check.Value, "value")
		}
	}
}

type parseTagsCase struct {
	Input         []string
	ExpectedMap   map[string]interface{}
	ExpectedError bool
	ExpectedMsg   string
}

func TestParseTags(t *testing.T) {
	t.Parallel()

	cases := []parseTagsCase{
		{[]string{"foo=bar", "true=false"}, map[string]interface{}{"foo": "bar", "true": "false"}, false, ""},
		{[]string{"foo=", "true=false"}, nil, true, `foo= is invalid, Must be in format key=value
Examples: rack=1 evironment=production`},
	}

	for _, check := range cases {
		actual, err := parseTags(check.Input)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoErr(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "map")
		}
	}
}

func TestTagsList(t *testing.T) {
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
			"memory": {},
			"cpu": {},
			"tags": {
				"warp": "8",
				"ncc": "1701"
			},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TagsList("enterprise")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), `=== enterprise Tags
ncc      1701
warp     8
`, "output")
}

func TestTagsSet(t *testing.T) {
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
				Tags: map[string]interface{}{
					"true": "false",
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {},
			"cpu": {},
			"tags": {
				"true": "false"
			},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TagsSet("foo", []string{"true=false"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying tags... done

=== foo Tags
true     false
`, "output")
}

func TestTagsUnset(t *testing.T) {
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
				Tags: map[string]interface{}{
					"ncc": nil,
				},
			}, r)
		}

		fmt.Fprintf(w, `{
			"owner": "jkirk",
			"app": "foo",
			"values": {},
			"memory": {},
			"cpu": {},
			"tags": {
				"warp": 8
			},
			"registry": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.TagsUnset("foo", []string{"ncc"})
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `Applying tags... done

=== foo Tags
warp     8
`, "output")
}
