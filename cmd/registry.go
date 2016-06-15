package cmd

import (
	"fmt"
	"strings"

	"github.com/deis/pkg/prettyprint"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// RegistryList lists an app's registry information.
func RegistryList(appID string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(c, appID)
	if checkAPICompatibility(c, err) != nil {
		return err
	}

	fmt.Printf("=== %s Registry\n", appID)

	registryMap := make(map[string]string)

	for key, value := range config.Registry {
		registryMap[key] = fmt.Sprintf("%v", value)
	}

	fmt.Print(prettyprint.PrettyTabs(registryMap, 5))

	return nil
}

// RegistrySet sets an app's registry information.
func RegistrySet(appID string, item []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	registryMap := parseInfos(item)

	fmt.Print("Applying registry information... ")

	quit := progress()
	configObj := api.Config{}
	configObj.Registry = registryMap

	_, err = config.Set(c, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(c, err) != nil {
		return err
	}

	fmt.Print("done\n\n")

	return RegistryList(appID)
}

// RegistryUnset removes an app's registry information.
func RegistryUnset(appID string, items []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Applying registry information... ")

	quit := progress()

	configObj := api.Config{}

	registryMap := make(map[string]interface{})

	for _, key := range items {
		registryMap[key] = nil
	}

	configObj.Registry = registryMap

	_, err = config.Set(c, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(c, err) != nil {
		return err
	}

	fmt.Print("done\n\n")

	return RegistryList(appID)
}

func parseInfos(items []string) map[string]interface{} {
	registryMap := make(map[string]interface{})

	for _, item := range items {
		key, value, err := parseInfo(item)

		if err != nil {
			fmt.Println(err)
			continue
		}

		registryMap[key] = value
	}

	return registryMap
}

func parseInfo(item string) (string, string, error) {
	parts := strings.SplitN(item, "=", 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf(`%s is invalid. Must be in format key=value
Examples: username=bob password=s3cur3pw1`, item)
	}

	if parts[0] != "username" && parts[0] != "password" {
		return "", "", fmt.Errorf(`%s is invalid. Valid keys are "username" or "password".`, parts[0])
	}

	return parts[0], parts[1], nil
}
