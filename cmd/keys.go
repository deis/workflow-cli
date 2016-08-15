package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/keys"
	"github.com/deis/workflow-cli/pkg/ssh"
	"github.com/deis/workflow-cli/settings"
)

// KeysList lists a user's keys.
func (d DeisCmd) KeysList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	keys, count, err := keys.List(s.Client, results)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Keys%s", s.Username, limitCount(len(keys), count))

	for _, key := range keys {
		d.Printf("%s %s...%s\n", key.ID, key.Public[:16], key.Public[len(key.Public)-10:])
	}
	return nil
}

// KeyRemove removes keys.
func (d DeisCmd) KeyRemove(keyID string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Printf("Removing %s SSH Key...", keyID)

	if err = keys.Delete(s.Client, keyID); checkAPICompatibility(s.Client, err, d.WErr) != nil {
		d.Println()
		return err
	}

	d.Println(" done")
	return nil
}

// KeyAdd adds keys.
func (d DeisCmd) KeyAdd(keyLocation string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	var key api.KeyCreateRequest

	if keyLocation == "" {
		ks, err := listKeys(d.WOut)
		if err != nil {
			return err
		}
		key, err = chooseKey(ks, os.Stdin, d.WOut)
		if err != nil {
			return err
		}
	} else {
		key, err = getKey(keyLocation)
		if err != nil {
			return err
		}
	}

	d.Printf("Uploading %s to deis...", filepath.Base(key.Name))

	if _, err = keys.New(s.Client, key.ID, key.Public); checkAPICompatibility(s.Client, err, d.WErr) != nil {
		d.Println()
		return err
	}

	d.Println(" done")
	return nil
}

func chooseKey(keys []api.KeyCreateRequest, input io.Reader,
	wOut io.Writer) (api.KeyCreateRequest, error) {
	fmt.Fprintln(wOut, "Found the following SSH public keys:")

	for i, key := range keys {
		fmt.Fprintf(wOut, "%d) %s %s\n", i+1, filepath.Base(key.Name), key.ID)
	}

	fmt.Fprintln(wOut, "0) Enter path to pubfile (or use keys:add <key_path>)")

	var selected string

	fmt.Fprint(wOut, "Which would you like to use with Deis? ")
	fmt.Fscanln(input, &selected)

	numSelected, err := strconv.Atoi(selected)

	if err != nil {
		return api.KeyCreateRequest{}, fmt.Errorf("%s is not a valid integer", selected)
	}

	if numSelected < 0 || numSelected > len(keys) {
		return api.KeyCreateRequest{}, fmt.Errorf("%d is not a valid option", numSelected)
	}

	if numSelected == 0 {
		var filename string

		fmt.Fprint(wOut, "Enter the path to the pubkey file: ")
		fmt.Scanln(&filename)

		return getKey(filename)
	}

	return keys[numSelected-1], nil
}

func listKeys(wOut io.Writer) ([]api.KeyCreateRequest, error) {
	folder := filepath.Join(settings.FindHome(), ".ssh")
	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	var keys []api.KeyCreateRequest

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".pub" {
			key, err := getKey(filepath.Join(folder, file.Name()))

			if err == nil {
				keys = append(keys, key)
			} else {
				fmt.Fprintln(wOut, err)
			}
		}
	}

	return keys, nil
}

func getKey(filename string) (api.KeyCreateRequest, error) {
	keyContents, err := ioutil.ReadFile(filename)

	if err != nil {
		return api.KeyCreateRequest{}, err
	}

	backupID := strings.Split(filepath.Base(filename), ".")[0]
	keyInfo, err := ssh.ParsePubKey(backupID, keyContents)
	if err != nil {
		return api.KeyCreateRequest{}, fmt.Errorf("%s is not a valid ssh key", filename)
	}
	return api.KeyCreateRequest{ID: keyInfo.ID, Public: keyInfo.Public, Name: filename}, nil
}
