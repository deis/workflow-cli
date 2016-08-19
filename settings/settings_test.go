package settings

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-cli/version"
)

const sFile string = `{"username":"t","ssl_verify":false,"controller":"http://foo.bar","token":"a"}`

func createTempProfile(contents string) (string, error) {
	name, err := ioutil.TempDir("", "client")

	if err != nil {
		return "", err
	}

	file := filepath.Join(name, "test.json")

	return file, ioutil.WriteFile(file, []byte(contents), 0775)
}

type comparison struct {
	key      interface{}
	expected interface{}
}

func TestLoadSave(t *testing.T) {
	t.Parallel()
	// Load profile from file and confirm it is correctly parsed.
	file, err := createTempProfile(sFile)
	if err != nil {
		t.Fatal(err)
	}

	s, err := Load(file)

	if err != nil {
		t.Fatal(err)
	}

	tests := []comparison{
		{
			key:      false,
			expected: s.Client.VerifySSL,
		},
		{
			key:      "a",
			expected: s.Client.Token,
		},
		{
			key:      "t",
			expected: s.Username,
		},
		{
			key:      "http://foo.bar",
			expected: s.Client.ControllerURL.String(),
		},
		{
			key:      100,
			expected: s.Limit,
		},
		{
			key:      "Deis Client " + version.Version,
			expected: s.Client.UserAgent,
		},
	}

	if err := checkComparisons(tests); err != nil {
		t.Error(err)
	}

	// Modify profile and confirm it is correctly saved
	s.Client.VerifySSL = true
	s.Client.Token = "b"
	s.Username = "c"
	s.Limit = 10

	u, err := url.Parse("http://deis.test")

	if err != nil {
		t.Fatal(err)
	}

	s.Client.ControllerURL = u

	if _, err = s.Save(file); err != nil {
		t.Fatal(err)
	}

	s, err = Load(file)

	if err != nil {
		t.Fatal(err)
	}

	tests = []comparison{
		{
			key:      true,
			expected: s.Client.VerifySSL,
		},
		{
			key:      "b",
			expected: s.Client.Token,
		},
		{
			key:      "c",
			expected: s.Username,
		},
		{
			key:      "http://deis.test",
			expected: s.Client.ControllerURL.String(),
		},
		{
			key:      10,
			expected: s.Limit,
		},
		{
			key:      "Deis Client " + version.Version,
			expected: s.Client.UserAgent,
		},
	}

	if err := checkComparisons(tests); err != nil {
		t.Error(err)
	}
}

func checkComparisons(tests []comparison) error {
	for _, check := range tests {
		if check.key != check.expected {
			return fmt.Errorf("Expected %v, Got %v", check.key, check.expected)
		}
	}

	return nil
}

func TestDeleteSettings(t *testing.T) {
	t.Parallel()

	file, err := createTempProfile("")
	if err != nil {
		t.Fatal(err)
	}

	if err := Delete(file); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(file); err == nil {
		t.Errorf("File %s exists, supposed to have been deleted.", file)
	}

	// Test the deleting an nonexistent settings file isn't an error.
	if err := Delete(file); err != nil {
		t.Fatal(err)
	}
}

func TestNotLoggedIn(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "client")

	if err != nil {
		t.Fatal(err)
	}

	_, err = Load(filepath.Join(name, "test.json"))
	assert.ExistsErr(t, err, "load error")
	if !strings.Contains(err.Error(), "Client configuration file not found") {
		t.Error("Expected configuration error, Got:", err.Error())
	}
}
