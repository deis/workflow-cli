package cmd

import (
	"fmt"
	"strings"

	"github.com/deis/pkg/prettyprint"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// TagsList lists an app's tags.
func (d DeisCmd) TagsList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Tags\n", appID)

	tagMap := make(map[string]string)

	for key, value := range config.Tags {
		tagMap[key] = fmt.Sprintf("%v", value)
	}

	d.Print(prettyprint.PrettyTabs(tagMap, 5))

	return nil
}

// TagsSet sets an app's tags.
func (d DeisCmd) TagsSet(appID string, tags []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	tagsMap, err := parseTags(tags)
	if err != nil {
		return err
	}

	d.Print("Applying tags... ")

	quit := progress(d.WOut)
	configObj := api.Config{}
	configObj.Tags = tagsMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TagsList(appID)
}

// TagsUnset removes an app's tags.
func (d DeisCmd) TagsUnset(appID string, tags []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Applying tags... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	tagsMap := make(map[string]interface{})

	for _, tag := range tags {
		tagsMap[tag] = nil
	}

	configObj.Tags = tagsMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TagsList(appID)
}

func parseTags(tags []string) (map[string]interface{}, error) {
	tagMap := make(map[string]interface{})

	for _, tag := range tags {
		key, value, err := parseTag(tag)

		if err != nil {
			return nil, err
		}

		tagMap[key] = value
	}

	return tagMap, nil
}

func parseTag(tag string) (string, string, error) {
	parts := strings.Split(tag, "=")

	if len(parts) != 2 {
		return "", "", fmt.Errorf(`%s is invalid, Must be in format key=value
Examples: rack=1 evironment=production`, tag)
	}

	return parts[0], parts[1], nil
}
