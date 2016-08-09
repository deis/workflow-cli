package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestScaleFail(t *testing.T) {
	t.Parallel()

	tmpDir, err := ioutil.TempDir("", "tmpdir")
	if err != nil {
		t.Fatalf("error creating temp directory (%s)", err)
	}

	file := filepath.Join(tmpDir, "client.json")

	data := []byte(`{"username":"test","ssl_verify":false,"controller":"http://deis.127.0.0.1.nip.io","token":"test","response_limit":0}`)
	if err := ioutil.WriteFile(file, data, 0644); err != nil {
		t.Fatal("error creating config file")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("failed to remove creds file from %s (%s)", tmpDir, err)
		}
	}()

	expected := "'web=-1' does not match the pattern 'type=num', ex: web=2\n"
	d := DeisCmd{ConfigFile: file}
	actual := d.PsScale("testApp", []string{"web=-1"})
	if actual.Error() != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

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
