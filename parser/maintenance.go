package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Maintenance displays all relevant commands for `deis maintenance`.
func Maintenance(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for maintenance:

maintenance:info   view maintenance mode of an application
maintenance:on     turn on maintenance for an app
maintenance:off    turn off maintenance for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "maintenance:info":
		return maintenanceInfo(argv, cmdr)
	case "maintenance:on":
		return maintenanceEnable(argv, cmdr)
	case "maintenance:off":
		return maintenanceDisable(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "maintenance" {
			argv[0] = "maintenance:info"
			return maintenanceInfo(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func maintenanceInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about the current application's maintenance state.

Usage: deis maintenance:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.MaintenanceInfo(safeGetValue(args, "--app"))
}

func maintenanceEnable(argv []string, cmdr cmd.Commander) error {
	usage := `
Enables maintenance mode for an app.

Usage: deis maintenance:on [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.MaintenanceEnable(safeGetValue(args, "--app"))
}

func maintenanceDisable(argv []string, cmdr cmd.Commander) error {
	usage := `
Disables maintenance mode for an app.

Usage: deis maintenance:off [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.MaintenanceDisable(safeGetValue(args, "--app"))
}
