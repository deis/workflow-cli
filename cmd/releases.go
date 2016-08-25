package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/deis/controller-sdk-go/releases"
)

// ReleasesList lists an app's releases.
func (d DeisCmd) ReleasesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	releases, count, err := releases.List(s.Client, appID, results)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Releases%s", appID, limitCount(len(releases), count))

	w := new(tabwriter.Writer)

	w.Init(d.WOut, 0, 8, 1, '\t', 0)
	for _, r := range releases {
		fmt.Fprintf(w, "v%d\t%s\t%s\n", r.Version, r.Created, r.Summary)
	}
	w.Flush()
	return nil
}

// ReleasesInfo prints info about a specific release.
func (d DeisCmd) ReleasesInfo(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	r, err := releases.Get(s.Client, appID, version)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Release v%d\n", appID, version)
	if r.Build != "" {
		d.Println("build:   ", r.Build)
	}
	d.Println("config:  ", r.Config)
	d.Println("owner:   ", r.Owner)
	d.Println("created: ", r.Created)
	d.Println("summary: ", r.Summary)
	d.Println("updated: ", r.Updated)
	d.Println("uuid:    ", r.UUID)

	return nil
}

// ReleasesRollback rolls an app back to a previous release.
func (d DeisCmd) ReleasesRollback(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if version == -1 {
		d.Print("Rolling back one release... ")
	} else {
		d.Printf("Rolling back to v%d... ", version)
	}

	quit := progress(d.WOut)
	newVersion, err := releases.Rollback(s.Client, appID, version)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("done, v%d\n", newVersion)

	return nil
}
