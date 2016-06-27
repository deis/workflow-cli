package cmd

import (
	"fmt"
	"regexp"

	"github.com/deis/pkg/prettyprint"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// LimitsList lists an app's limits.
func LimitsList(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Limits\n\n", appID)

	fmt.Println("--- Memory")
	if len(config.Memory) == 0 {
		fmt.Println("Unlimited")
	} else {
		memoryMap := make(map[string]string)

		for key, value := range config.Memory {
			memoryMap[key] = fmt.Sprintf("%v", value)
		}

		fmt.Print(prettyprint.PrettyTabs(memoryMap, 5))
	}

	fmt.Println("\n--- CPU")
	if len(config.CPU) == 0 {
		fmt.Println("Unlimited")
	} else {
		cpuMap := make(map[string]string)

		for key, value := range config.CPU {
			cpuMap[key] = value.(string)
		}

		fmt.Print(prettyprint.PrettyTabs(cpuMap, 5))
	}

	return nil
}

// LimitsSet sets an app's limits.
func LimitsSet(appID string, limits []string, limitType string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	limitsMap, err := parseLimits(limits)
	if err != nil {
		return err
	}

	fmt.Print("Applying limits... ")

	quit := progress()
	configObj := api.Config{}

	if limitType == "cpu" {
		configObj.CPU = limitsMap
	} else {
		configObj.Memory = limitsMap
	}

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Print("done\n\n")

	return LimitsList(appID)
}

// LimitsUnset removes an app's limits.
func LimitsUnset(appID string, limits []string, limitType string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Applying limits... ")

	quit := progress()

	configObj := api.Config{}

	valuesMap := make(map[string]interface{})

	for _, limit := range limits {
		valuesMap[limit] = nil
	}

	if limitType == "cpu" {
		configObj.CPU = valuesMap
	} else {
		configObj.Memory = valuesMap
	}

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Print("done\n\n")

	return LimitsList(appID)
}

func parseLimits(limits []string) (map[string]interface{}, error) {
	limitsMap := make(map[string]interface{})

	for _, limit := range limits {
		key, value, err := parseLimit(limit)

		if err != nil {
			return nil, fmt.Errorf("Error parsing limit %s: %s", limit, err)
		}

		limitsMap[key] = value
	}

	return limitsMap, nil
}

func parseLimit(limit string) (string, string, error) {
	regex := regexp.MustCompile("^([A-z]+)=([0-9]+[bkmgBKMG]{1,2}|[0-9.]{1,5}|[0-9.]{1,5}[m]{0,1})$")

	if !regex.MatchString(limit) {
		return "", "", fmt.Errorf(`%s doesn't fit format type=#unit or type=#
Examples: web=2G worker=500M web=300`, limit)
	}

	capture := regex.FindStringSubmatch(limit)

	return capture[1], capture[2], nil
}
