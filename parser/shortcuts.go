package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Shortcuts displays all relevant shortcuts for the CLI.
func Shortcuts(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for shortcuts:

shortcuts:list       list all relevant shortcuts for the CLI

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "shortcuts:list":
		return shortcutsList(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "shortcuts" {
			argv[0] = "shortcuts:list"
			return shortcutsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func shortcutsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists all relevant shortcuts for the CLI

Usage: deis shortcuts:list
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.ShortcutsList()
}
