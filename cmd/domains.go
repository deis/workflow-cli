package cmd

import "github.com/deis/controller-sdk-go/domains"

// DomainsList lists domains registered with an app.
func (d DeisCmd) DomainsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	domains, count, err := domains.List(s.Client, appID, results)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Domains%s", appID, limitCount(len(domains), count))

	for _, domain := range domains {
		d.Println(domain.Domain)
	}
	return nil
}

// DomainsAdd adds a domain to an app.
func (d DeisCmd) DomainsAdd(appID, domain string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Adding %s to %s... ", domain, appID)

	quit := progress(d.WOut)
	_, err = domains.New(s.Client, appID, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// DomainsRemove removes a domain registered with an app.
func (d DeisCmd) DomainsRemove(appID, domain string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s... ", domain, appID)

	quit := progress(d.WOut)
	err = domains.Delete(s.Client, appID, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Println("done")
	return nil
}
