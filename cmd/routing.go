package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/appsettings"
)

// RoutingInfo provides information about the status of app routing.
func (d DeisCmd) RoutingInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if *appSettings.Routable {
		fmt.Println("Routing is enabled.")
	} else {
		fmt.Println("Routing is disabled.")
	}
	return nil
}

// RoutingEnable enables an app from being exposed by the router.
func (d DeisCmd) RoutingEnable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	fmt.Printf("Enabling routing for %s... ", appID)

	quit := progress()
	appSettings := api.AppSettings{Routable: api.NewRoutable()}
	_, err = appsettings.Set(s.Client, appID, appSettings)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}

// RoutingDisable disables an app from being exposed by the router.
func (d DeisCmd) RoutingDisable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	fmt.Printf("Disabling routing for %s... ", appID)

	quit := progress()
	appSettings := api.AppSettings{Routable: api.NewRoutable()}
	*appSettings.Routable = false
	_, err = appsettings.Set(s.Client, appID, appSettings)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}
