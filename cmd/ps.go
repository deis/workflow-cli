package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/ps"
)

// PsList lists an app's processes.
func (d DeisCmd) PsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	processes, _, err := ps.List(s.Client, appID, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	printProcesses(appID, processes)

	return nil
}

// PsScale scales an app's processes.
func (d DeisCmd) PsScale(appID string, targets []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	targetMap := make(map[string]int)
	regex := regexp.MustCompile("^([a-z0-9]+)=([0-9]+)$")

	for _, target := range targets {
		if regex.MatchString(target) {
			captures := regex.FindStringSubmatch(target)
			targetMap[captures[1]], err = strconv.Atoi(captures[2])

			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("'%s' does not match the pattern 'type=num', ex: web=2\n", target)
		}
	}

	fmt.Printf("Scaling processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress()

	err = ps.Scale(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	processes, _, err := ps.List(s.Client, appID, s.Limit)
	if err != nil {
		return err
	}

	printProcesses(appID, processes)
	return nil
}

// PsRestart restarts an app's processes.
func (d DeisCmd) PsRestart(appID, target string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	psType, psName := "", ""
	if target != "" {
		psType, psName = parseType(target, appID)
	}

	fmt.Printf("Restarting processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress()

	processes, err := ps.Restart(s.Client, appID, psType, psName)
	quit <- true
	<-quit
	if err == deis.ErrPodNotFound {
		return fmt.Errorf("Could not find process type %s in app %s", psType, appID)
	} else if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(processes) == 0 {
		fmt.Println("Could not find any processes to restart")
	} else {
		fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
		printProcesses(appID, processes)
	}

	return nil
}

func printProcesses(appID string, processes []api.Pods) {
	psMap := ps.ByType(processes)

	fmt.Printf("=== %s Processes\n", appID)

	for psType, procs := range psMap {
		fmt.Printf("--- %s:\n", psType)

		for _, proc := range procs {
			fmt.Printf("%s %s (%s)\n", proc.Name, proc.State, proc.Release)
		}
	}
}

func parseType(target string, appID string) (string, string) {
	psType, psName := "", ""

	if strings.Contains(target, "-") {
		replaced := strings.Replace(target, appID+"-", "", 1)
		parts := strings.Split(replaced, "-")
		// the API requires the type, for now
		// regex matches against how Deployment pod name is constructed
		regex := regexp.MustCompile("[0-9]{8,10}-[a-z0-9]{5}$")
		if regex.MatchString(replaced) {
			psType = parts[0]
		} else {
			psType = parts[1]
		}
		// process name is the full pod
		psName = target
	} else {
		psType = target
	}

	return psType, psName
}
