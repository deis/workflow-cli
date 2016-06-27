package cmd

import (
	"fmt"

	"github.com/deis/controller-sdk-go/domains"
)

// DomainsList lists domains registered with an app.
func DomainsList(appID string, results int) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	domains, count, err := domains.List(s.Client, appID, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Domains%s", appID, limitCount(len(domains), count))

	for _, domain := range domains {
		fmt.Println(domain.Domain)
	}
	return nil
}

// DomainsAdd adds a domain to an app.
func DomainsAdd(appID, domain string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Adding %s to %s... ", domain, appID)

	quit := progress()
	_, err = domains.New(s.Client, appID, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

// DomainsRemove removes a domain registered with an app.
func DomainsRemove(appID, domain string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Removing %s from %s... ", domain, appID)

	quit := progress()
	err = domains.Delete(s.Client, appID, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Println("done")
	return nil
}
