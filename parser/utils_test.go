package parser

import (
	"bytes"
	"testing"

	"github.com/deis/workflow-cli/pkg/testutil"
)

func TestSafeGet(t *testing.T) {
	t.Parallel()

	expected := "foo"

	test := make(map[string]interface{}, 1)
	test["test"] = "foo"

	actual := safeGetValue(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGetNil(t *testing.T) {
	t.Parallel()

	expected := ""

	test := make(map[string]interface{}, 1)
	test["test"] = nil

	actual := safeGetValue(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGetInt(t *testing.T) {
	t.Parallel()

	expected := 1

	test := make(map[string]interface{}, 1)
	test["test"] = "1"

	actual := safeGetInt(test, "test")

	if expected != actual {
		t.Errorf("Expected %d, Got %d", expected, actual)
	}

	if actual = safeGetInt(test, "foo"); actual != 0 {
		t.Errorf("Expected 0, Got %d", actual)
	}
}

func TestPrintHelp(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDeisCmd{WOut: &b, ConfigFile: cf}

	usage := ""

	if !printHelp([]string{"ps", "--help"}, usage, cmdr) {
		t.Error("Expected true")
	}

	if !printHelp([]string{"ps", "-h"}, usage, cmdr) {
		t.Error("Expected true")
	}

	if printHelp([]string{"ps"}, usage, cmdr) {
		t.Error("Expected false")
	}

	if printHelp([]string{"ps", "--foo"}, usage, cmdr) {
		t.Error("Expected false")
	}
}
