package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Perms routes perms commands to their specific function.
func Perms(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for perms:

perms:list            list permissions granted on an app
perms:create          create a new permission for a user
perms:delete          delete a permission for a user

Use 'deis help perms:[command]' to learn more.
`

	switch argv[0] {
	case "perms:list":
		return permsList(argv, cmdr)
	case "perms:create":
		return permCreate(argv, cmdr)
	case "perms:delete":
		return permDelete(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "perms" {
			argv[0] = "perms:list"
			return permsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func permsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists all users with permission to use an app, or lists all users with system
administrator privileges.

Usage: deis perms:list [-a --app=<app>|--admin|--admin --limit=<num>]

Options:
  -a --app=<app>
    lists all users with permission to <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    lists all users with system administrator privileges.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	admin := args["--admin"].(bool)

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.PermsList(app, admin, results)
}

func permCreate(argv []string, cmdr cmd.Commander) error {
	usage := `
Gives another user permission to use an app, or gives another user
system administrator privileges.

Usage: deis perms:create <username> [-a --app=<app>|--admin]

Arguments:
  <username>
    the name of the new user.

Options:
  -a --app=<app>
    grants <username> permission to use <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    grants <username> system administrator privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	username := args["<username>"].(string)
	admin := args["--admin"].(bool)

	return cmdr.PermCreate(app, username, admin)
}

func permDelete(argv []string, cmdr cmd.Commander) error {
	usage := `
Revokes another user's permission to use an app, or revokes another user's system
administrator privileges.

Usage: deis perms:delete <username> [-a --app=<app>|--admin]

Arguments:
  <username>
    the name of the user.

Options:
  -a --app=<app>
    revokes <username> permission to use <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    revokes <username> system administrator privileges.`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	username := args["<username>"].(string)
	admin := args["--admin"].(bool)

	return cmdr.PermDelete(app, username, admin)
}
