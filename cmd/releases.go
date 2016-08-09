package cmd

import (
	"fmt"
	"os"
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
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Releases%s", appID, limitCount(len(releases), count))

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
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
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Release v%d\n", appID, version)
	if r.Build != "" {
		fmt.Println("build:   ", r.Build)
	}
	fmt.Println("config:  ", r.Config)
	fmt.Println("owner:   ", r.Owner)
	fmt.Println("created: ", r.Created)
	fmt.Println("summary: ", r.Summary)
	fmt.Println("updated: ", r.Updated)
	fmt.Println("uuid:    ", r.UUID)

	return nil
}

// ReleasesRollback rolls an app back to a previous release.
func (d DeisCmd) ReleasesRollback(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if version == -1 {
		fmt.Print("Rolling back one release... ")
	} else {
		fmt.Printf("Rolling back to v%d... ", version)
	}

	quit := progress()
	newVersion, err := releases.Rollback(s.Client, appID, version)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("done, v%d\n", newVersion)

	return nil
}
