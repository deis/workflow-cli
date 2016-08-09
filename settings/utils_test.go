package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

type confgCases struct {
	Input    string
	Expected string
}

func TestSelectSettings(t *testing.T) {
	t.Parallel()
	cases := []confgCases{
		{"test", filepath.Join(FindHome(), ".deis", "test.json")},
		{"", filepath.Join(FindHome(), ".deis", "client.json")},
		{"~/test.json", "~/test.json"},
		{"/opt/test.json", "/opt/test.json"},
	}

	for _, check := range cases {
		assert.Equal(t, locateSettingsFile(check.Input), check.Expected, "case")
	}

	// Check that env variable is used.
	location := "/test/test.json"
	os.Setenv("DEIS_PROFILE", location)
	assert.Equal(t, locateSettingsFile(""), location, "case")
}
