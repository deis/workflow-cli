package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/version"
)

// DefaultResponseLimit is the default number of responses to return on requests that can
// be limited.
const DefaultResponseLimit = 100

// UserAgent is the user agent used by the CLI
var UserAgent = "Deis Client v" + version.Version

type settingsFile struct {
	Username   string `json:"username"`
	VerifySSL  bool   `json:"ssl_verify"`
	Controller string `json:"controller"`
	Token      string `json:"token"`
	Limit      int    `json:"response_limit"`
}

// Settings is the settings object created from the settings file.
type Settings struct {
	Username string
	Limit    int
	Client   *deis.Client
}

// Load loads a new client from a settings file.
func Load(cf string) (*Settings, error) {
	filename := locateSettingsFile(cf)

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(`Client configuration file not found at: %s
Are you logged in? Use 'deis login' or 'deis register' to get started.`, filename)
		}

		return nil, err
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	sF := settingsFile{}
	if err = json.Unmarshal(contents, &sF); err != nil {
		return nil, err
	}

	c, err := deis.New(sF.VerifySSL, sF.Controller, sF.Token)

	if err != nil {
		return nil, err
	}

	// Set a custom user agent
	c.UserAgent = UserAgent

	settings := Settings{}
	settings.Username = sF.Username
	settings.Client = c

	// If users have defined a custom response limit, respect it.
	if sF.Limit > 0 {
		settings.Limit = sF.Limit
	} else {
		settings.Limit = DefaultResponseLimit
	}

	return &settings, nil
}

// Save settings to a file
func (s *Settings) Save(cf string) (string, error) {
	settings := settingsFile{Username: s.Username, VerifySSL: s.Client.VerifySSL,
		Controller: s.Client.ControllerURL.String(), Token: s.Client.Token, Limit: s.Limit}

	settingsContents, err := json.Marshal(settings)

	if err != nil {
		return "", err
	}

	if err = os.MkdirAll(filepath.Join(FindHome(), "/.deis/"), 0775); err != nil {
		return "", err
	}

	filename := locateSettingsFile(cf)

	return filename, ioutil.WriteFile(filename, settingsContents, 0775)
}

// Delete user's settings file.
func Delete(cf string) error {
	filename := locateSettingsFile(cf)

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
