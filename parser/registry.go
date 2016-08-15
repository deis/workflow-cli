package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Registry routes registry commands to their specific function
func Registry(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for registry:

registry:list        list registry info for an app
registry:set         set registry info for an app
registry:unset       unset registry info for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "registry:list":
		return registryList(argv, cmdr)
	case "registry:set":
		return registrySet(argv, cmdr)
	case "registry:unset":
		return registryUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "registry" {
			argv[0] = "registry:list"
			return registryList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func registryList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists registry information for an application.

Usage: deis registry:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.RegistryList(safeGetValue(args, "--app"))
}

func registrySet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets registry information for an application. These credentials are the same as those used for
'docker login' to the private registry.

Usage: deis registry:set [options] <key>=<value>...

Arguments:
  <key>
    the uniquely identifiable name for logging into the registry. Valid keys are "username" or
    "password"
  <value>
    the value of said environment variable. For example, "bob" or "mysecretpassword"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	info := args["<key>=<value>"].([]string)

	return cmdr.RegistrySet(app, info)
}

func registryUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets registry information for an application.

Usage: deis registry:unset [options] <key>...

Arguments:
  <key> the registry key to unset, for example: "username" or "password"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	key := args["<key>"].([]string)

	return cmdr.RegistryUnset(app, key)
}
