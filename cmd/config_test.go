package cmd

import (
	"testing"

	"github.com/arschles/assert"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	_, err := parseConfig([]string{"FOO=bar", "CAR star"})
	assert.ExistsErr(t, err, "config")

	actual, err := parseConfig([]string{"FOO=bar"})
	assert.NoErr(t, err)
	assert.Equal(t, actual, map[string]interface{}{"FOO": "bar"}, "map")
}
