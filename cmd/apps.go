package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deis/pkg/prettyprint"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/apps"
	"github.com/deis/controller-sdk-go/config"
	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/pkg/webbrowser"
	"github.com/deis/workflow-cli/settings"
)

// AppCreate creates an app.
func AppCreate(id string, buildpack string, remote string, noRemote bool) error {
	s, err := settings.Load()
	if err != nil {
		return err
	}

	fmt.Print("Creating Application... ")
	quit := progress()
	app, err := apps.New(s.Client, id)

	quit <- true
	<-quit

	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("done, created %s\n", app.ID)

	if buildpack != "" {
		configValues := api.Config{
			Values: map[string]interface{}{
				"BUILDPACK_URL": buildpack,
			},
		}
		if _, err = config.Set(s.Client, app.ID, configValues); checkAPICompatibility(s.Client, err) != nil {
			return err
		}
	}

	if !noRemote {
		if err = git.CreateRemote(s.Client.ControllerURL.Host, remote, app.ID); err != nil {
			if err.Error() == "exit status 128" {
				fmt.Println("To replace the existing git remote entry, run:")
				fmt.Printf("  git remote rename deis deis.old && deis git:remote -a %s\n", app.ID)
			}
			return err
		}
	}

	if noRemote {
		fmt.Printf("If you want to add a git remote for this app later, use `deis git:remote -a %s`\n", app.ID)
	} else {
		fmt.Println("remote available at", git.RemoteURL(s.Client.ControllerURL.Host, app.ID))
	}

	return nil
}

// AppsList lists apps on the Deis controller.
func AppsList(results int) error {
	s, err := settings.Load()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	apps, count, err := apps.List(s.Client, results)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== Apps%s", limitCount(len(apps), count))

	for _, app := range apps {
		fmt.Println(app.ID)
	}
	return nil
}

// AppInfo prints info about app.
func AppInfo(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("=== %s Application\n", app.ID)
	fmt.Println("updated: ", app.Updated)
	fmt.Println("uuid:    ", app.UUID)
	fmt.Println("created: ", app.Created)
	fmt.Println("url:     ", app.URL)
	fmt.Println("owner:   ", app.Owner)
	fmt.Println("id:      ", app.ID)

	fmt.Println()
	// print the app processes
	if err = PsList(app.ID, defaultLimit); err != nil {
		return err
	}

	fmt.Println()
	// print the app domains
	if err = DomainsList(app.ID, defaultLimit); err != nil {
		return err
	}

	fmt.Println()

	return nil
}

// AppOpen opens an app in the default webbrowser.
func AppOpen(appID string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(s.Client, appID)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	u := app.URL
	if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
		u = "http://" + u
	}

	return webbrowser.Webbrowser(u)
}

// AppLogs returns the logs from an app.
func AppLogs(appID string, lines int) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	logs, err := apps.Logs(s.Client, appID, lines)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	return printLogs(logs)
}

// printLogs prints each log line with a color matched to its category.
func printLogs(logs string) error {
	for _, log := range strings.Split(logs, `\n`) {
		category := "unknown"
		parts := strings.Split(strings.Split(log, " -- ")[0], " ")
		category = parts[0]
		colorVars := map[string]string{
			"Color": chooseColor(category),
			"Log":   log,
		}
		fmt.Println(prettyprint.ColorizeVars("{{.V.Color}}{{.V.Log}}{{.C.Default}}", colorVars))
	}

	return nil
}

// AppRun runs a one time command in the app.
func AppRun(appID, command string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Running '%s'...\n", command)

	out, err := apps.Run(s.Client, appID, command)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if out.ReturnCode == 0 {
		fmt.Print(out.Output)
	} else {
		fmt.Fprint(os.Stderr, out.Output)
	}

	os.Exit(out.ReturnCode)
	return nil
}

// AppDestroy destroys an app.
func AppDestroy(appID, confirm string) error {
	gitSession := false

	s, err := settings.Load()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = git.DetectAppName(s.Client.ControllerURL.Host)

		if err != nil {
			return err
		}

		gitSession = true
	}

	if confirm == "" {
		fmt.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, appID, appID, appID)

		fmt.Scanln(&confirm)
	}

	if confirm != appID {
		return fmt.Errorf("App %s does not match confirm %s, aborting.", appID, confirm)
	}

	startTime := time.Now()
	fmt.Printf("Destroying %s...\n", appID)

	if err = apps.Delete(s.Client, appID); checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	if gitSession {
		return git.DeleteRemote(appID)
	}

	return nil
}

// AppTransfer transfers app ownership to another user.
func AppTransfer(appID, username string) error {
	s, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Transferring %s to %s... ", appID, username)

	err = apps.Transfer(s.Client, appID, username)
	if checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	fmt.Println("done")

	return nil
}
