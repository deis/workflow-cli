package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Domains routes domain commands to their specific function.
func Domains(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for domains:

domains:add           bind a domain to an application
domains:list          list domains bound to an application
domains:remove        unbind a domain from an application

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "domains:add":
		return domainsAdd(argv, cmdr)
	case "domains:list":
		return domainsList(argv, cmdr)
	case "domains:remove":
		return domainsRemove(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "domains" {
			argv[0] = "domains:list"
			return domainsList(argv, cmdr)
		}

		PrintUsage()
		return nil
	}
}

func domainsAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Binds a domain to an application.

Usage: deis domains:add <domain> [options]

Arguments:
  <domain>
    the domain name to be bound to the application, such as 'domain.deisapp.com'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	domain := safeGetValue(args, "<domain>")

	return cmdr.DomainsAdd(app, domain)
}

func domainsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists domains bound to an application.

Usage: deis domains:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")

	return cmdr.DomainsList(app, results)
}

func domainsRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
Unbinds a domain for an application.

Usage: deis domains:remove <domain> [options]

Arguments:
  <domain>
    the domain name to be removed from the application.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	domain := safeGetValue(args, "<domain>")

	return cmdr.DomainsRemove(app, domain)
}
