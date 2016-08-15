package cmd

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/deis/controller-sdk-go/builds"
)

// BuildsList lists an app's builds.
func (d DeisCmd) BuildsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	builds, count, err := builds.List(s.Client, appID, results)
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Printf("=== %s Builds%s", appID, limitCount(len(builds), count))

	for _, build := range builds {
		d.Println(build.UUID, build.Created)
	}
	return nil
}

// BuildsCreate creates a build for an app.
func (d DeisCmd) BuildsCreate(appID, image, procfile string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	procfileMap := make(map[string]string)

	if procfile != "" {
		if procfileMap, err = parseProcfile([]byte(procfile)); err != nil {
			return err
		}
	} else if _, err := os.Stat("Procfile"); err == nil {
		contents, err := ioutil.ReadFile("Procfile")
		if err != nil {
			return err
		}

		if procfileMap, err = parseProcfile(contents); err != nil {
			return err
		}
	}

	d.Print("Creating build... ")
	quit := progress(d.WOut)
	_, err = builds.New(s.Client, appID, image, procfileMap)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err, d.WErr) != nil {
		return err
	}

	d.Println("done")

	return nil
}

func parseProcfile(procfile []byte) (map[string]string, error) {
	procfileMap := make(map[string]string)
	return procfileMap, yaml.Unmarshal(procfile, &procfileMap)
}
