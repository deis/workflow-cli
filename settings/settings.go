package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	client "github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/version"
)

// UserAgent is the user agent used by the CLI
const UserAgent = "Deis Client v" + version.Version

type settingsFile struct {
	Username   string `json:"username"`
	VerifySSL  bool   `json:"ssl_verify"`
	Controller string `json:"controller"`
	Token      string `json:"token"`
	Limit      int    `json:"response_limit"`
}

// Load loads a new client from a settings file.
func Load() (*client.Client, error) {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Not logged in. Use 'deis login' or 'deis register' to get started.")
		}

		return nil, err
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	settings := settingsFile{}
	if err = json.Unmarshal(contents, &settings); err != nil {
		return nil, err
	}

	c, err := client.New(settings.VerifySSL, settings.Controller, settings.Token, settings.Username)

	if err != nil {
		return nil, err
	}

	// If users have defined a custom response limit, respect it.
	if settings.Limit > 0 {
		c.ResponseLimit = settings.Limit
	}

	// Set a custom user agent
	c.UserAgent = UserAgent

	return c, nil
}

// Save settings to a file
func Save(c *client.Client) error {
	settings := settingsFile{Username: c.Username, VerifySSL: c.VerifySSL,
		Controller: c.ControllerURL.String(), Token: c.Token, Limit: c.ResponseLimit}

	settingsContents, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Join(FindHome(), "/.deis/"), 0775); err != nil {
		return err
	}

	return ioutil.WriteFile(locateSettingsFile(), settingsContents, 0775)
}

// Delete user's settings file.
func Delete() error {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}
