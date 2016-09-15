package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestPermsListUsers(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/perms/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
			"users": [
                "baz",
                "bar"
            ]
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", false, -1)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `=== foo's Users
baz
bar
`, "output")
}

func TestPermsListUsersLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/apps/foo/perms/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
			"users": [
                "baz"
            ]
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", false, 1)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `=== foo's Users
baz
`, "output")
}

func TestPermsListAdmins(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
            "count": 2,
            "next": null,
            "previous": null,
			"results": [
                {
                    "username": "fred",
                    "is_superuser": true
                },
                {
                    "username": "bob",
                    "is_superuser": true
                }
            ]
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", true, -1)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `=== Administrators
fred
bob
`, "output")
}

func TestPermsListAdminsLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/admin/perms/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)

		fmt.Fprintf(w, `{
            "count": 2,
            "next": null,
            "previous": null,
			"results": [
                {
                    "username": "fred",
                    "is_superuser": true
                }
            ]
		}`)
	})

	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.PermsList("foo", true, 1)
	assert.NoErr(t, err)

	assert.Equal(t, testutil.StripProgress(b.String()), `=== Administrators (1 of 2)
fred
`, "output")
}
