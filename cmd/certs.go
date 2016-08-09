package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/certs"
	"github.com/deis/workflow-cli/settings"
)

// CertsList lists certs registered with the controller.
func (d DeisCmd) CertsList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	certList, _, err := certs.List(s.Client, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(certList) == 0 {
		fmt.Println("No certs")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderLine(true)
	table.SetHeader([]string{"Name", "Common Name", "SubjectAltName", "Expires", "Fingerprint", "Domains", "Updated", "Created"})
	for _, cert := range certList {
		domains := strings.Join(cert.Domains[:], ",")
		san := strings.Join(cert.SubjectAltName[:], ",")

		// Make dates more readable
		now := time.Now()
		expires := cert.Expires.Time.Format("2 Jan 2006")
		created := cert.Created.Time.Format("2 Jan 2006")
		updated := cert.Updated.Time.Format("2 Jan 2006")

		if cert.Expires.Time.Before(now) {
			expires += " (expired)"
		} else {
			// Ghetto solution
			expires += " (in"
			year := cert.Expires.Time.Year() - now.Year()
			month := cert.Expires.Time.Month() - now.Month()
			day := cert.Expires.Time.Day() - now.Day()

			if year > 0 {
				expires += fmt.Sprintf(" %d year", year)
				if year > 1 {
					expires += "s"
				}
			} else if month > 0 {
				expires += fmt.Sprintf(" %d month", month)
				if month > 1 {
					expires += "s"
				}
			} else if day != 0 {
				// special handling on negative days
				if day < 0 {
					day *= -1
				}

				expires += fmt.Sprintf(" %d day", day)
				if day > 1 {
					expires += "s"
				}
			}
			expires += ")"
		}

		// show a shorter version of the fingerprint
		fingerprint := cert.Fingerprint[:5] + "[...]" + cert.Fingerprint[len(cert.Fingerprint)-5:]

		table.Append([]string{cert.Name, cert.CommonName, san, expires, fingerprint, domains, updated, created})
	}
	table.Render()

	return nil
}

// CertAdd adds a cert to the controller.
func (d DeisCmd) CertAdd(cert string, key string, name string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	fmt.Print("Adding SSL endpoint... ")
	quit := progress()
	err = doCertAdd(s.Client, cert, key, name)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

func doCertAdd(c *deis.Client, cert string, key string, name string) error {
	certFile, err := ioutil.ReadFile(cert)
	if err != nil {
		return err
	}

	keyFile, err := ioutil.ReadFile(key)
	if err != nil {
		return err
	}

	_, err = certs.New(c, string(certFile), string(keyFile), name)
	return checkAPICompatibility(c, err)
}

// CertRemove deletes a cert from the controller.
func (d DeisCmd) CertRemove(name string) error {
	s, err := settings.Load(d.ConfigFile)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("Removing %s... ", name)
	quit := progress()

	err = certs.Delete(s.Client, name)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

// CertInfo gets info about certficiate
func (d DeisCmd) CertInfo(name string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	cert, err := certs.Get(s.Client, name)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	domains := strings.Join(cert.Domains[:], ",")
	if domains == "" {
		domains = "No connected domains"
	}

	san := strings.Join(cert.SubjectAltName[:], ",")
	if san == "" {
		san = "N/A"
	}

	fmt.Printf("=== %s Certificate\n", cert.Name)
	fmt.Println("Common Name(s):    ", cert.CommonName)
	fmt.Println("Expires At:        ", cert.Expires)
	fmt.Println("Starts At:         ", cert.Starts)
	fmt.Println("Fingerprint:       ", cert.Fingerprint)
	fmt.Println("Subject Alt Name:  ", san)
	fmt.Println("Issuer:            ", cert.Issuer)
	fmt.Println("Subject:           ", cert.Subject)
	fmt.Println()
	fmt.Println("Connected Domains: ", domains)
	fmt.Println("Owner:             ", cert.Owner)
	fmt.Println("Created:           ", cert.Created)
	fmt.Println("Updated:           ", cert.Updated)

	return nil
}

// CertAttach attaches a certificate to a domain
func (d DeisCmd) CertAttach(name, domain string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	fmt.Printf("Attaching certificate %s to domain %s... ", name, domain)
	quit := progress()

	err = certs.Attach(s.Client, name, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) == nil {
		fmt.Println("done")
	}

	return err
}

// CertDetach detaches a certificate from a domain
func (d DeisCmd) CertDetach(name, domain string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	fmt.Printf("Detaching certificate %s from domain %s... ", name, domain)
	quit := progress()

	err = certs.Detach(s.Client, name, domain)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Println("done")
	return nil
}
