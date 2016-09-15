package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Whitelist displays all relevant commands for `deis whitelist`.
func Whitelist(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for whitelist:

whitelist:add           adds addresses to the application's whitelist
whitelist:list          list addresses in the application's whitelist
whitelist:remove        remove addresses from the application's whitelist

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "whitelist:add":
		return whitelistAdd(argv, cmdr)
	case "whitelist:list":
		return whitelistList(argv, cmdr)
	case "whitelist:remove":
		return whitelistRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "whitelist" {
			argv[0] = "whitelist:list"
			return whitelistList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func whitelistAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Adds addresses to an application whitelist.

Usage: deis whitelist:add <addresses> [options]

Arguments:
  <addresses>
    comma-delimited list of addresses(using IP or CIDR notation) to be whitelisted for the application, such as '1.2.3.4' or '1.2.3.4,0.0.0.0/0'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	addresses := safeGetValue(args, "<addresses>")

	return cmdr.WhitelistAdd(app, addresses)
}

func whitelistList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists whitelisted addresses for an application.

Usage: deis whitelist:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.WhitelistList(app)
}

func whitelistRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Removes addresses from an application whitelist.

Usage: deis whitelist:remove <addresses> [options]

Arguments:
  <addresses>
    comma-delimited list of addresses(using IP or CIDR notation) to be whitelisted for the application, such as '1.2.3.4' or "1.2.3.4,0.0.0.0/0".

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	addresses := safeGetValue(args, "<addresses>")

	return cmdr.WhitelistRemove(app, addresses)
}
