package cmd

import (
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/appsettings"
)

// AutoscaleList tells the informations about app's autoscale status
func (d *DeisCmd) AutoscaleList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Autoscale\n\n", appID)

	if appSettings.Autoscale == nil {
		d.Println("No autoscale rules found.")
	} else {
		for process, kv := range appSettings.Autoscale {
			d.Println("--- " + process + ":")
			d.Println(*kv)
		}
	}

	return nil
}

// AutoscaleSet sets autoscale options for the app.
func (d *DeisCmd) AutoscaleSet(appID string, processType string, min int, max int, CPUPercent int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Applying autoscale settings for process type %s on %s... ", processType, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		processType: {
			Min:        min,
			Max:        max,
			CPUPercent: CPUPercent,
		},
	}
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Autoscale: data})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AutoscaleUnset removes autoscale for the app.
func (d *DeisCmd) AutoscaleUnset(appID string, processType string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing autoscale for process type %s on %s... ", processType, appID)

	quit := progress(d.WOut)
	data := map[string]*api.Autoscale{
		processType: nil,
	}
	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Autoscale: data})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}
