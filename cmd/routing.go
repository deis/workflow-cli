package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

func RoutingInfo(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if config.Routable == true {
		fmt.Println("Routing is enabled.")
	} else {
		fmt.Println("Routing is disabled.")
	}
	return nil
}

// RoutingEnable enables an app from being exposed by the router.
func RoutingEnable(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Enabling routing for %s... ", appID)

	quit := progress()
	_, err = config.Set(s.Client, appID, api.NewConfig())

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}

// RoutingDisable disables an app from being exposed by the router.
func RoutingDisable(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Disabling routing for %s... ", appID)

	quit := progress()
	configObj := api.NewConfig()
	configObj.Routable = false
	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")
	return nil
}
