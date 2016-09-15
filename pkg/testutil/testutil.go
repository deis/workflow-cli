package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/settings"
)

// TestServer represents a test HTTP server along with a path to a config file
type TestServer struct {
	Server *httptest.Server
	Mux    *http.ServeMux
}

// NewTestServer sets up a test HTTP Server without a configuration file.
func NewTestServer() *TestServer {
	// test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &TestServer{
		Server: server,
		Mux:    mux,
	}
}

// Close closes the test HTTP server.
func (t *TestServer) Close() {
	t.Server.Close()
}

// NewTestServerAndClient sets up a test HTTP Server along with a configuration file to talk to it
func NewTestServerAndClient() (string, *TestServer, error) {
	server := NewTestServer()

	name, err := ioutil.TempDir("", "client")
	if err != nil {
		server.Close()
		return "", nil, err
	}

	filename := filepath.Join(name, "test.json")

	client, err := deis.New(false, server.Server.URL, "")
	if err != nil {
		server.Close()
		return "", nil, err
	}

	config := settings.Settings{
		Username: "test",
		Client:   client,
	}

	filename, err = config.Save(filename)
	if err != nil {
		server.Close()
		return "", nil, err
	}
	return filename, server, nil
}

// AssertBody asserts the value of the body of a request.
func AssertBody(t *testing.T, expected interface{}, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	value, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != string(value) {
		t.Errorf("Expected body '%s' actually got '%s'\n", value, body)
	}
}

// StripProgress strips the output from the progress method
func StripProgress(input string) string {
	first := strings.Index(input, "\b")
	// If \b charecter not part of string
	if first == -1 {
		return input
	}
	last := strings.LastIndex(input, "\b")

	// return string without \b or the characters it deletes.
	return input[:first-(last-first+1)] + input[last+1:]
}

// SetHeaders sets standard headers for requests
func SetHeaders(w http.ResponseWriter) {
	w.Header().Add("DEIS_API_VERSION", deis.APIVersion)
}
