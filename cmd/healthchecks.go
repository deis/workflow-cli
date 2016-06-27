package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// HealthchecksList lists an app's healthchecks.
func HealthchecksList(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Healthchecks\n\n", appID)

	fmt.Println("--- Liveness")
	if livenessProbe, found := config.Healthcheck["livenessProbe"]; found {
		fmt.Println(livenessProbe)
	} else {
		fmt.Println("No liveness probe configured.")
	}

	fmt.Println("\n--- Readiness")
	if readinessProbe, found := config.Healthcheck["readinessProbe"]; found {
		fmt.Println(readinessProbe)
	} else {
		fmt.Println("No readiness probe configured.")
	}
	return nil
}

// HealthchecksSet sets an app's healthchecks.
func HealthchecksSet(appID, healthcheckType string, probe *api.Healthcheck) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Applying %s healthcheck... ", healthcheckType)

	quit := progress()
	configObj := api.Config{}
	configObj.Healthcheck = make(map[string]*api.Healthcheck)

	configObj.Healthcheck[healthcheckType] = probe

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return HealthchecksList(appID)
}

// HealthchecksUnset removes an app's healthchecks.
func HealthchecksUnset(appID string, healthchecks []string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Removing healthchecks... ")

	quit := progress()

	configObj := api.Config{}

	healthcheckMap := make(map[string]*api.Healthcheck)

	for _, healthcheck := range healthchecks {
		healthcheckMap[healthcheck] = nil
	}

	configObj.Healthcheck = healthcheckMap

	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return HealthchecksList(appID)
}
