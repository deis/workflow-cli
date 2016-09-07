package cmd

import "github.com/deis/controller-sdk-go/tls"

// TLSInfo prints info about the TLS settings for the given app.
func (d *DeisCmd) TLSInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	tls, err := tls.Info(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s TLS\n", appID)
	d.Println(tls)

	return nil
}

// TLSEnable enables the router to enforce https-only requests to the application.
func (d *DeisCmd) TLSEnable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Enabling https-only requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.Enable(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// TLSDisable disables the router to enforce https-only requests to the application.
func (d *DeisCmd) TLSDisable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Disabling https-only requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.Disable(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
