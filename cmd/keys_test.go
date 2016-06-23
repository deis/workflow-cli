package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/deis/controller-sdk-go/api"
)

func TestGetKey(t *testing.T) {
	t.Parallel()

	file, err := ioutil.TempFile("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc test@example.com")

	expected := api.KeyCreateRequest{
		ID:     "test@example.com",
		Public: string(toWrite),
		Name:   file.Name(),
	}

	if _, err = file.Write(toWrite); err != nil {
		t.Fatal(err)
	}

	key, err := getKey(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, key) {
		t.Errorf("Expected %v, Got %v", expected, key)
	}
}

func TestGetKeyNoComment(t *testing.T) {
	t.Parallel()

	file, err := ioutil.TempFile("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc")

	expected := api.KeyCreateRequest{
		ID:     filepath.Base(file.Name()),
		Public: string(toWrite),
		Name:   file.Name(),
	}

	if _, err = file.Write(toWrite); err != nil {
		t.Fatal(err)
	}

	key, err := getKey(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, key) {
		t.Errorf("Expected %v, Got %v", expected, key)
	}
}

func TestListKeys(t *testing.T) {
	name, err := ioutil.TempDir("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("HOME", name)

	folder := filepath.Join(name, ".ssh")

	if err = os.Mkdir(folder, 0755); err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc test@example.com")
	fileNames := []string{"test1.pub", "test2.pub"}

	expected := []api.KeyCreateRequest{
		{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   filepath.Join(folder, fileNames[0]),
		},
		{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   filepath.Join(folder, fileNames[1]),
		},
	}

	for _, file := range fileNames {
		if err = ioutil.WriteFile(filepath.Join(folder, file), toWrite, 0775); err != nil {
			t.Fatal(err)
		}
	}

	keys, err := listKeys()

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, keys) {
		t.Errorf("Expected %v, Got %v", expected, keys)
	}
}

type chooseKeyCases struct {
	Reader   io.Reader
	Expected string
}

func TestChooseKey(t *testing.T) {
	testKeys := []api.KeyCreateRequest{
		api.KeyCreateRequest{},
	}

	checks := []chooseKeyCases{
		chooseKeyCases{
			Reader:   strings.NewReader("-1"),
			Expected: "-1 is not a valid option",
		},
		chooseKeyCases{
			Reader:   strings.NewReader("2"),
			Expected: "2 is not a valid option",
		},
		chooseKeyCases{
			Reader:   strings.NewReader("a"),
			Expected: "a is not a valid integer",
		},
	}

	for _, check := range checks {
		_, err := chooseKey(testKeys, check.Reader)

		if err.Error() != check.Expected {
			t.Errorf("Expected %s, Got %s", check.Expected, err.Error())
		}
	}
}
