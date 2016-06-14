package client

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

const sFile string = `{"username":"t","ssl_verify":false,"controller":"http://d.t","token":"a","response_limit": 50}`

func createTempProfile(contents string) error {
	name, err := ioutil.TempDir("", "client")

	if err != nil {
		return err
	}

	os.Unsetenv("DEIS_PROFILE")
	os.Setenv("HOME", name)
	folder := filepath.Join(name, ".deis")
	if err = os.Mkdir(folder, 0755); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(folder, "client.json"), []byte(contents), 0775); err != nil {
		return err
	}

	return nil
}

type comparison struct {
	key      interface{}
	expected interface{}
}

func TestLoadSave(t *testing.T) {
	// Load profile from file and confirm it is correctly parsed.
	if err := createTempProfile(sFile); err != nil {
		t.Fatal(err)
	}

	client, err := New()

	if err != nil {
		t.Fatal(err)
	}

	tests := []comparison{
		comparison{
			key:      false,
			expected: client.SSLVerify,
		},
		comparison{
			key:      "a",
			expected: client.Token,
		},
		comparison{
			key:      "t",
			expected: client.Username,
		},
		comparison{
			key:      "http://d.t",
			expected: client.ControllerURL.String(),
		},
		comparison{
			key:      50,
			expected: client.ResponseLimit,
		},
	}

	checkComparisons(tests, t)

	// Modify profile and confirm it is correctly saved
	client.SSLVerify = true
	client.Token = "b"
	client.Username = "c"
	client.ResponseLimit = 100

	u, err := url.Parse("http://deis.test")

	if err != nil {
		t.Fatal(err)
	}

	client.ControllerURL = *u

	if err = client.Save(); err != nil {
		t.Fatal(err)
	}

	client, err = New()

	if err != nil {
		t.Fatal(err)
	}

	tests = []comparison{
		comparison{
			key:      true,
			expected: client.SSLVerify,
		},
		comparison{
			key:      "b",
			expected: client.Token,
		},
		comparison{
			key:      "c",
			expected: client.Username,
		},
		comparison{
			key:      "http://deis.test",
			expected: client.ControllerURL.String(),
		},
		comparison{
			key:      100,
			expected: client.ResponseLimit,
		},
	}

	checkComparisons(tests, t)
}

func checkComparisons(tests []comparison, t *testing.T) {
	for _, check := range tests {
		if check.key != check.expected {
			t.Errorf("Expected %v, Got %v", check.key, check.expected)
		}
	}
}

func TestDeleteSettings(t *testing.T) {
	if err := createTempProfile(""); err != nil {
		t.Fatal(err)
	}

	if err := Delete(); err != nil {
		t.Fatal(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Errorf("File %s exists, supposed to have been deleted.", file)
	}
}
