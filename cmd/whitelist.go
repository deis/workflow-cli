package cmd

import (
	"strings"

	"github.com/deis/controller-sdk-go/whitelist"
)

// WhitelistList lists the addresses whitelisted for app
func (d *DeisCmd) WhitelistList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	whitelist, err := whitelist.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Whitelisted Addresses", appID)

	for _, ip := range whitelist.Addresses {
		d.Println(ip)
	}
	return nil
}

// WhitelistAdd adds the addresses to the app's Whitelist.
func (d *DeisCmd) WhitelistAdd(appID, IPs string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Adding %s to %s whitelist...\n", IPs, appID)

	quit := progress(d.WOut)
	_, err = whitelist.Add(s.Client, appID, strings.Split(IPs, ","))
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// WhitelistRemove deletes the addresses from the app's Whitelist.
func (d *DeisCmd) WhitelistRemove(appID, IPs string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s whitelist...\n", IPs, appID)

	quit := progress(d.WOut)
	err = whitelist.Delete(s.Client, appID, strings.Split(IPs, ","))
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
