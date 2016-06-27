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
func KeysList(results int) error {
	s, err := settings.Load()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	keys, count, err := keys.List(s.Client, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Keys%s", s.Username, limitCount(len(keys), count))

	for _, key := range keys {
		fmt.Printf("%s %s...%s\n", key.ID, key.Public[:16], key.Public[len(key.Public)-10:])
	}
	return nil
}

// KeyRemove removes keys.
func KeyRemove(keyID string) error {
	s, err := settings.Load()

	if err != nil {
		return err
	}

	fmt.Printf("Removing %s SSH Key...", keyID)

	if err = keys.Delete(s.Client, keyID); checkAPICompatibility(s.Client, err) != nil {
		fmt.Println()
		return err
	}

	fmt.Println(" done")
	return nil
}

// KeyAdd adds keys.
func KeyAdd(keyLocation string) error {
	s, err := settings.Load()

	if err != nil {
		return err
	}

	var key api.KeyCreateRequest

	if keyLocation == "" {
		ks, err := listKeys()
		if err != nil {
			return err
		}
		key, err = chooseKey(ks, os.Stdin)
		if err != nil {
			return err
		}
	} else {
		key, err = getKey(keyLocation)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Uploading %s to deis...", filepath.Base(key.Name))

	if _, err = keys.New(s.Client, key.ID, key.Public); checkAPICompatibility(s.Client, err) != nil {
		fmt.Println()
		return err
	}

	fmt.Println(" done")
	return nil
}

func chooseKey(keys []api.KeyCreateRequest, input io.Reader) (api.KeyCreateRequest, error) {
	fmt.Println("Found the following SSH public keys:")

	for i, key := range keys {
		fmt.Printf("%d) %s %s\n", i+1, filepath.Base(key.Name), key.ID)
	}

	fmt.Println("0) Enter path to pubfile (or use keys:add <key_path>)")

	var selected string

	fmt.Print("Which would you like to use with Deis? ")
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

		fmt.Print("Enter the path to the pubkey file: ")
		fmt.Scanln(&filename)

		return getKey(filename)
	}

	return keys[numSelected-1], nil
}

func listKeys() ([]api.KeyCreateRequest, error) {
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
				fmt.Println(err)
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
