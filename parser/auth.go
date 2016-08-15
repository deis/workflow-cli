package parser

import (
	"fmt"

	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Auth routes auth commands to the specific function.
func Auth(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for auth:

auth:register          register a new user
auth:login             authenticate against a controller
auth:logout            clear the current user session
auth:passwd            change the password for the current user
auth:whoami            display the current user
auth:cancel            remove the current user account
auth:regenerate        regenerate user tokens

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "auth:register":
		return authRegister(argv, cmdr)
	case "auth:login":
		return authLogin(argv, cmdr)
	case "auth:logout":
		return authLogout(argv, cmdr)
	case "auth:passwd":
		return authPasswd(argv, cmdr)
	case "auth:whoami":
		return authWhoami(argv, cmdr)
	case "auth:cancel":
		return authCancel(argv, cmdr)
	case "auth:regenerate":
		return authRegenerate(argv, cmdr)
	case "auth":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func authRegister(argv []string, cmdr cmd.Commander) error {
	usage := `
Registers a new user with a Deis controller.

Usage: deis auth:register <controller> [options]

Arguments:
  <controller>
    fully-qualified controller URI, e.g. 'http://deis.local3.deisapp.com/'

Options:
  --username=<username>
    provide a username for the new account.
  --password=<password>
    provide a password for the new account.
  --email=<email>
    provide an email address.
  --ssl-verify=false
    disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	email := safeGetValue(args, "--email")
	sslVerify := false

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "true" {
		sslVerify = true
	}

	return cmdr.Register(controller, username, password, email, sslVerify)
}

func authLogin(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs in by authenticating against a controller.

Usage: deis auth:login <controller> [options]

Arguments:
  <controller>
    a fully-qualified controller URI, e.g. "http://deis.local3.deisapp.com/".

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --ssl-verify=false
    disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	sslVerify := false

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "true" {
		sslVerify = true
	}

	return cmdr.Login(controller, username, password, sslVerify)
}

func authLogout(argv []string, cmdr cmd.Commander) error {
	usage := `
Logs out from a controller and clears the user session.

Usage: deis auth:logout

Options:
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmdr.Logout()
}

func authPasswd(argv []string, cmdr cmd.Commander) error {
	usage := `
Changes the password for the current user.

Usage: deis auth:passwd [options]

Options:
  --password=<password>
    the current password for the account.
  --new-password=<new-password>
    the new password for the account.
  --username=<username>
    the account's username.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	newPassword := safeGetValue(args, "--new-password")

	return cmdr.Passwd(username, password, newPassword)
}

func authWhoami(argv []string, cmdr cmd.Commander) error {
	usage := `
Displays the currently logged in user.

Usage: deis auth:whoami [options]

Options:
  --all
    fetch a more detailed description about the user.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.Whoami(args["--all"].(bool))
}

func authCancel(argv []string, cmdr cmd.Commander) error {
	usage := `
Cancels and removes the current account.

Usage: deis auth:cancel [options]

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --yes
    force "yes" when prompted.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	yes := args["--yes"].(bool)

	return cmdr.Cancel(username, password, yes)
}

func authRegenerate(argv []string, cmdr cmd.Commander) error {
	usage := `
Regenerates auth token, defaults to regenerating token for the current user.

Usage: deis auth:regenerate [options]

Options:
  -u --username=<username>
    specify user to regenerate. Requires admin privileges.
  --all
    regenerate token for every user. Requires admin privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	all := args["--all"].(bool)

	return cmdr.Regenerate(username, all)
}
