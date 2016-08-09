package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/appsettings"
)

// MaintenanceInfo tells the informations about app's maintenance status
func (d DeisCmd) MaintenanceInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if *appSettings.Maintenance {
		fmt.Println("Maintenance mode is on.")
	} else {
		fmt.Println("Maintenance mode is off.")
	}
	return nil
}

// MaintenanceEnable turns on the maintenance for the app.
func (d DeisCmd) MaintenanceEnable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	fmt.Printf("Enabling maintenance mode for %s... ", appID)

	quit := progress()
	b := true
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Maintenance: &b})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}

// MaintenanceDisable turns off the maintenance for the app.
func (d DeisCmd) MaintenanceDisable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	fmt.Printf("Disabling maintenance mode for %s... ", appID)

	quit := progress()
	b := false
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Maintenance: &b})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}
