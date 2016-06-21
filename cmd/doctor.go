package cmd

import (
	"fmt"
	"strings"

	"github.com/deis/controller-sdk-go/doctor"
	"github.com/deis/workflow-cli/settings"
)

// DoctorSend sends Dcotor info to deis servers.
func DoctorSend() error {
	dc, err := settings.Load()
	url := strings.SplitAfterN(dc.ControllerURL.Host, ".", 2)[1]
	dcHost := "deis-workflow-manager." + url
	dc.ControllerURL.Host = dcHost
	uidurl, err := doctor.Send(dc)
	if err != nil {
		return err
	}

	fmt.Printf("the information received by deis doctor can be seen at the following url https://doctor-staging.deis.com/%s \n", uidurl)
	return nil
}
