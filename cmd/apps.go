package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/apps"
	"github.com/deis/controller-sdk-go/config"
	"github.com/deis/controller-sdk-go/domains"
	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/pkg/logging"
	"github.com/deis/workflow-cli/pkg/webbrowser"
	"github.com/deis/workflow-cli/settings"
)

// AppCreate creates an app.
func (d *DeisCmd) AppCreate(id, buildpack, remote string, noRemote bool) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Print("Creating Application... ")
	quit := progress(d.WOut)
	app, err := apps.New(s.Client, id)

	quit <- true
	<-quit

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, created %s\n", app.ID)

	if buildpack != "" {
		configValues := api.Config{
			Values: map[string]interface{}{
				"BUILDPACK_URL": buildpack,
			},
		}
		if _, err = config.Set(s.Client, app.ID, configValues); d.checkAPICompatibility(s.Client, err) != nil {
			return err
		}
	}

	if !noRemote {
		if err = git.CreateRemote(git.DefaultCmd, s.Client.ControllerURL.Host, remote, app.ID); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("fatal: remote %s already exists.", remote)) {
				msg := "A git remote with the name %s already exists. To overwrite this remote run:\n"
				msg += "deis git:remote --force --remote %s --app %s"
				return fmt.Errorf(msg, remote, remote, app.ID)
			}
			return err
		}

		d.Printf(remoteCreationMsg, remote, app.ID)
	}

	if noRemote {
		d.Printf("If you want to add a git remote for this app later, use `deis git:remote -a %s`\n", app.ID)
	}

	return nil
}

// AppsList lists apps on the Deis controller.
func (d *DeisCmd) AppsList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	apps, count, err := apps.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== Apps%s", limitCount(len(apps), count))

	for _, app := range apps {
		d.Println(app.ID)
	}
	return nil
}

// AppInfo prints info about app.
func (d *DeisCmd) AppInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	url, err := d.appURL(s, appID)
	if err != nil {
		return err
	}

	if url == "" {
		url = fmt.Sprintf(noDomainAssignedMsg, appID)
	}

	d.Printf("=== %s Application\n", app.ID)
	d.Println("updated: ", app.Updated)
	d.Println("uuid:    ", app.UUID)
	d.Println("created: ", app.Created)
	d.Println("url:     ", url)
	d.Println("owner:   ", app.Owner)
	d.Println("id:      ", app.ID)

	d.Println()
	// print the app processes
	if err = d.PsList(app.ID, defaultLimit); err != nil {
		return err
	}

	d.Println()
	// print the app domains
	if err = d.DomainsList(app.ID, defaultLimit); err != nil {
		return err
	}

	d.Println()

	return nil
}

// AppOpen opens an app in the default webbrowser.
func (d *DeisCmd) AppOpen(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	u, err := d.appURL(s, appID)
	if err != nil {
		return err
	}

	if u == "" {
		return fmt.Errorf(noDomainAssignedMsg, appID)
	}

	if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
		u = "http://" + u
	}

	return webbrowser.Webbrowser(u)
}

// AppLogs returns the logs from an app.
func (d *DeisCmd) AppLogs(appID string, lines int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	logs, err := apps.Logs(s.Client, appID, lines)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	for _, log := range strings.Split(strings.TrimRight(logs, `\n`), `\n`) {
		logging.PrintLog(os.Stdout, log)
	}

	return nil
}

// AppRun runs a one time command in the app.
func (d *DeisCmd) AppRun(appID, command string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Running '%s'...\n", command)

	out, err := apps.Run(s.Client, appID, command)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if out.ReturnCode == 0 {
		d.Print(out.Output)
	} else {
		d.PrintErr(out.Output)
	}

	os.Exit(out.ReturnCode)
	return nil
}

// AppDestroy destroys an app.
func (d *DeisCmd) AppDestroy(appID, confirm string) error {
	gitSession := false

	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)

		if err != nil {
			return err
		}

		gitSession = true
	}

	if confirm == "" {
		d.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, appID, appID, appID)

		fmt.Scanln(&confirm)
	}

	if confirm != appID {
		return fmt.Errorf("App %s does not match confirm %s, aborting.", appID, confirm)
	}

	startTime := time.Now()
	d.Printf("Destroying %s...\n", appID)

	if err = apps.Delete(s.Client, appID); d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	if gitSession {
		return d.GitRemove(appID)
	}

	return nil
}

// AppTransfer transfers app ownership to another user.
func (d *DeisCmd) AppTransfer(appID, username string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Transferring %s to %s... ", appID, username)

	err = apps.Transfer(s.Client, appID, username)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

const noDomainAssignedMsg = "No domain assigned to %s"

// appURL grabs the first domain an app has and returns this.
func (d *DeisCmd) appURL(s *settings.Settings, appID string) (string, error) {
	domains, _, err := domains.List(s.Client, appID, 1)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return "", err
	}

	if len(domains) == 0 {
		return "", nil
	}

	return expandURL(s.Client.ControllerURL.Host, domains[0].Domain), nil
}

// expandURL expands an app url if necessary.
func expandURL(host, u string) string {
	if strings.Contains(u, ".") {
		// If domain is a full url.
		return u
	}

	// If domain is a subdomain, look up the controller url and replace the subdomain.
	parts := strings.Split(host, ".")
	parts[0] = u
	return strings.Join(parts, ".")
}
