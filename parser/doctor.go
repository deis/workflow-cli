package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Doctor routes doctor commands to their specific function.
func Doctor(argv []string) error {
	usage := `
Valid commands for Doctor:

doctor:send    sends doctor information

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "doctor:send":
		return doctorSend(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "doctor" {
			argv[0] = "doctor:send"
			return doctorSend(argv)
		}

		PrintUsage()
		return nil
	}
}

func doctorSend(argv []string) error {
	usage := `
Sends doctor information to deis servers

Usage: deis doctor:send [options]

`

	_, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}
	return cmd.DoctorSend()
}
