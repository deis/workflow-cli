package settings

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/deis/workflow-cli/version"
)

const sFile string = `{"username":"t","ssl_verify":false,"controller":"http://foo.bar","token":"a","response_limit": 50}`

func createTempProfile(contents string) error {
	name, err := ioutil.TempDir("", "client")

	if err != nil {
		return err
	}

	os.Unsetenv("DEIS_PROFILE")
	SetHome(name)
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

	c, err := Load()

	if err != nil {
		t.Fatal(err)
	}

	tests := []comparison{
		comparison{
			key:      false,
			expected: c.VerifySSL,
		},
		comparison{
			key:      "a",
			expected: c.Token,
		},
		comparison{
			key:      "t",
			expected: c.Username,
		},
		comparison{
			key:      "http://foo.bar",
			expected: c.ControllerURL.String(),
		},
		comparison{
			key:      50,
			expected: c.ResponseLimit,
		},
		comparison{
			key:      "Deis Client v" + version.Version,
			expected: c.UserAgent,
		},
	}

	if err := checkComparisons(tests); err != nil {
		t.Error(err)
	}

	// Modify profile and confirm it is correctly saved
	c.VerifySSL = true
	c.Token = "b"
	c.Username = "c"
	c.ResponseLimit = 100

	u, err := url.Parse("http://deis.test")

	if err != nil {
		t.Fatal(err)
	}

	c.ControllerURL = u

	if err = Save(c); err != nil {
		t.Fatal(err)
	}

	c, err = Load()

	if err != nil {
		t.Fatal(err)
	}

	tests = []comparison{
		comparison{
			key:      true,
			expected: c.VerifySSL,
		},
		comparison{
			key:      "b",
			expected: c.Token,
		},
		comparison{
			key:      "c",
			expected: c.Username,
		},
		comparison{
			key:      "http://deis.test",
			expected: c.ControllerURL.String(),
		},
		comparison{
			key:      100,
			expected: c.ResponseLimit,
		},
		comparison{
			key:      "Deis Client v" + version.Version,
			expected: c.UserAgent,
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
