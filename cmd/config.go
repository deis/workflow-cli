package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/deis/pkg/prettyprint"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// ConfigList lists an app's config.
func (d *DeisCmd) ConfigList(appID string, format string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	keys := sortKeys(config.Values)

	var configOutput *bytes.Buffer = new(bytes.Buffer)

	switch format {
	case "oneline":
		for i, key := range keys {
			sep := " "
			if i == len(keys)-1 {
				sep = "\n"
			}
			fmt.Fprintf(configOutput, "%s=%s%s", key, config.Values[key], sep)
		}
	case "diff":
		for _, key := range keys {
			fmt.Fprintf(configOutput, "%s=%s\n", key, config.Values[key])
		}
	default:
		fmt.Fprintf(configOutput, "=== %s Config\n", appID)

		configMap := make(map[string]string)

		// config.Values is type interface, so it needs to be converted to a string
		for _, key := range keys {
			configMap[key] = fmt.Sprintf("%v", config.Values[key])
		}

		fmt.Fprint(configOutput, prettyprint.PrettyTabs(configMap, 6))
	}

	d.Print(configOutput)
	return nil
}

// ConfigSet sets an app's config variables.
func (d *DeisCmd) ConfigSet(appID string, configVars []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	configMap, err := parseConfig(configVars)
	if err != nil {
		return err
	}

	if value, ok := configMap["SSH_KEY"]; ok {
		sshKey, err := parseSSHKey(value.(string))
		if err != nil {
			return err
		}
		configMap["SSH_KEY"] = base64.StdEncoding.EncodeToString([]byte(sshKey))
	}

	// NOTE(bacongobbler): check if the user is using the old way to set healthchecks. If so,
	// send them a deprecation notice.
	for key := range configMap {
		if strings.Contains(key, "HEALTHCHECK_") {
			d.Println(`Hey there! We've noticed that you're using 'deis config:set HEALTHCHECK_URL'
to set up healthchecks. This functionality has been deprecated. In the future, please use
'deis healthchecks' to set up application health checks. Thanks!`)
		}
	}

	d.Print("Creating config... ")

	quit := progress(d.WOut)
	configObj := api.Config{Values: configMap}
	configObj, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if release, ok := configObj.Values["WORKFLOW_RELEASE"]; ok {
		d.Printf("done, %s\n\n", release)
	} else {
		d.Print("done\n\n")
	}

	return d.ConfigList(appID, "")
}

// ConfigUnset removes a config variable from an app.
func (d *DeisCmd) ConfigUnset(appID string, configVars []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Removing config... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	valuesMap := make(map[string]interface{})

	for _, configVar := range configVars {
		valuesMap[configVar] = nil
	}

	configObj.Values = valuesMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.ConfigList(appID, "")
}

// ConfigPull pulls an app's config to a file.
func (d *DeisCmd) ConfigPull(appID string, interactive bool, overwrite bool) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	configVars, err := config.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	stat, err := os.Stdout.Stat()

	if err != nil {
		return err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		d.Print(formatConfig(configVars.Values))
		return nil
	}

	filename := ".env"

	if !overwrite {
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("%s already exists, pass -o to overwrite", filename)
		}
	}

	if interactive {
		contents, err := ioutil.ReadFile(filename)

		if err != nil {
			return err
		}
		localConfigVars := strings.Split(string(contents), "\n")

		configMap, err := parseConfig(localConfigVars[:len(localConfigVars)-1])
		if err != nil {
			return err
		}

		for key, value := range configVars.Values {
			localValue, ok := configMap[key]

			if ok {
				if value != localValue {
					var confirm string
					d.Printf("%s: overwrite %s with %s? (y/N) ", key, localValue, value)

					fmt.Scanln(&confirm)

					if strings.ToLower(confirm) == "y" {
						configMap[key] = value
					}
				}
			} else {
				configMap[key] = value
			}
		}

		return ioutil.WriteFile(filename, []byte(formatConfig(configMap)), 0755)
	}

	return ioutil.WriteFile(filename, []byte(formatConfig(configVars.Values)), 0755)
}

// ConfigPush pushes an app's config from a file.
func (d *DeisCmd) ConfigPush(appID, fileName string, verbose bool) error {
	stat, err := os.Stdin.Stat()

	if err != nil {
		return err
	}

	var contents []byte
	if ((stat.Mode()&os.ModeNamedPipe) != 0 || stat.Size() > 0) && fileName == ".no-file" {
		if verbose {
			d.Println("Getting input from stdin")
		}
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(os.Stdin)
		contents = buffer.Bytes()
	} else {
		if verbose {
			d.Println("Not getting input from stdin")
		}
		if fileName == ".no-file" {
			fileName = ".env"
		}
		contents, err = ioutil.ReadFile(fileName)

		if err != nil {
			return err
		}
	}

	file := strings.Split(string(contents), "\n")
	config := []string{}

	for _, configVar := range file {
		// If file has CRLF encoding, the default on windows, strip the CR
		configVar = strings.Trim(configVar, "\r")
		if len(configVar) > 0 {
			config = append(config, configVar)
		}
	}

	return d.ConfigSet(appID, config)
}

func parseConfig(configVars []string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*)=([\s\S]*)$`)
	for _, config := range configVars {
		// Skip config that starts with an comment
		if config[0] == '#' {
			continue
		}

		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			configMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test\n", config)
		}
	}

	return configMap, nil
}

func parseSSHKey(value string) (string, error) {
	sshRegex := regexp.MustCompile("^-----BEGIN (DSA|RSA|EC|OPENSSH) PRIVATE KEY-----")

	if sshRegex.MatchString(value) {
		return value, nil
	}

	// NOTE(felixbuenemann): check if the current value is already a base64 encoded key.
	// This is the case if it was fetched using "deis config:pull".
	contents, err := base64.StdEncoding.DecodeString(value)

	if err == nil && sshRegex.MatchString(string(contents)) {
		return string(contents), nil
	}

	// NOTE(felixbuenemann): check if the value is a path to a private key.
	if _, err := os.Stat(value); err == nil {
		contents, err := ioutil.ReadFile(value)

		if err != nil {
			return "", err
		}

		if sshRegex.MatchString(string(contents)) {
			return string(contents), nil
		}
	}

	return "", fmt.Errorf("Could not parse SSH private key:\n %s", value)
}

func formatConfig(configVars map[string]interface{}) string {
	var formattedConfig string

	keys := sortKeys(configVars)
	for _, key := range keys {
		formattedConfig += fmt.Sprintf("%s=%v\n", key, configVars[key])
	}

	return formattedConfig
}

func sortKeys(kv map[string]interface{}) []string {
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
