package testutil

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/arschles/assert"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// TestStripProgress ensures StripProgress strips what is expected.
func TestStripProgress(t *testing.T) {
	testInput := "Lorem ipsum dolar sit amet"
	expectedOutput := "Lorem ipsum dolar sit amet"

	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")

	testInput = "Lorem ipsum dolar sit amet...\b\b\b"
	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")
}

// TestAssertBody ensures AssertBody correctly marshals into the interface.
func TestAssertBody(t *testing.T) {

	b := nopCloser{bytes.NewBufferString(`{"data":{"lorem":"ipsum"},"dolar":["sit","amet"]}`)}

	sampleRequest := http.Request{
		Body: b,
	}

	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"lorem": "ipsum",
		},
		"dolar": []string{
			"sit",
			"amet",
		},
	}

	AssertBody(t, expected, &sampleRequest)
}
