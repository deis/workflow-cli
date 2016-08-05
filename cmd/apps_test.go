package cmd

import "testing"

type expandURLCases struct {
	Input    string
	Expected string
}

func TestExpandUrl(t *testing.T) {
	checks := []expandURLCases{
		{
			Input:    "test.com",
			Expected: "test.com",
		},
		{
			Input:    "test",
			Expected: "test.foo.com",
		},
	}

	for _, check := range checks {
		out := expandURL("deis.foo.com", check.Input)

		if out != check.Expected {
			t.Errorf("Expected %s, Got %s", check.Expected, out)
		}
	}
}
