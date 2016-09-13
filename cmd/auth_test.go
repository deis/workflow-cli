package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestRegister(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{}`)
	})

	username := "test-user"
	password := "test-password"
	email := "test-email@example.com"

	server.Mux.HandleFunc("/v2/auth/register/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "email": "`+email+`",
  "username": "`+username+`",
  "first_name": "",
  "last_name": "",
  "is_superuser": false,
  "is_staff": false,
  "groups": [],
  "user_permissions": [],
  "last_login": "2016-09-13T18:55:54Z",
  "date_joined": "2016-09-13T18:55:54Z",
  "is_active": true
}`)
	})

	server.Mux.HandleFunc("/v2/auth/login/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.Register(server.Server.URL, username, password, email, true)
	assert.NoErr(t, err)
	expected := fmt.Sprintf("Registered %s\nLogged in as %s\nConfiguration file written to %s\n", username, username, cf)

	assert.Equal(t, b.String(), expected, "output")
}

func TestLogin(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{}`)
	})

	server.Mux.HandleFunc("/v2/auth/login/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	username := "test-user"
	err = cmdr.Login(server.Server.URL, username, "test-pass", true)
	assert.NoErr(t, err)
	expected := fmt.Sprintf("Logged in as %s\nConfiguration file written to %s\n", username, cf)
	assert.Equal(t, b.String(), expected, "output")
}

func TestLogout(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.Logout()
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Logged out\n", "output")
}

func TestPasswd(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/auth/passwd/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	// Change own password.
	err = cmdr.Passwd("", "old-pass", "new-pass")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Password change succeeded.\n", "output")
	b.Reset()

	// Change another user's password.
	err = cmdr.Passwd("another-user", "old-pass", "new-pass")
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Password change succeeded.\n", "output")
}

func TestCancel(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{}`)
	})

	server.Mux.HandleFunc("/v2/auth/cancel/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	server.Mux.HandleFunc("/v2/auth/login/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	username := "test-user"
	err = cmdr.Cancel(username, "", true)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Account cancelled\n", "output")
}

func TestWhoami(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/auth/whoami/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "email": "test@example.com",
  "username": "test",
  "first_name": "",
  "last_name": "",
  "is_superuser": true,
  "is_staff": true,
  "groups": [],
  "user_permissions": [],
  "last_login": "2016-09-12T22:15:26Z",
  "date_joined": "2016-09-12T22:15:26Z",
  "is_active": true
}`)
	})

	err = cmdr.Whoami(false)
	assert.NoErr(t, err)
	expected := fmt.Sprintf("You are test at %s\n", server.Server.URL)
	assert.Equal(t, b.String(), expected, "output")
}

func TestRegenerate(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DeisCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/auth/tokens/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{}`)
	})

	err = cmdr.Regenerate("test", false)
	assert.NoErr(t, err)
	assert.Equal(t, b.String(), "Token Regenerated\n", "output")
}
