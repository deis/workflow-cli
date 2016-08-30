package parser

import (
	"time"

	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Certs routes certs commands to their specific function.
func Certs(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for certs:

certs:list            list SSL certificates for an app
certs:add             add an SSL certificate to an app
certs:remove          remove an SSL certificate from an app
certs:info            get detailed informaton about the certificate
certs:attach          attach an SSL certificate to a domain
certs:detach          detach an SSL certificate from a domain

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "certs:list":
		return certsList(argv, cmdr)
	case "certs:add":
		return certAdd(argv, cmdr)
	case "certs:remove":
		return certRemove(argv, cmdr)
	case "certs:info":
		return certInfo(argv, cmdr)
	case "certs:attach":
		return certAttach(argv, cmdr)
	case "certs:detach":
		return certDetach(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "certs" {
			argv[0] = "certs:list"
			return certsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func certsList(argv []string, cmdr cmd.Commander) error {
	usage := `
Show certificate information for an SSL application.

Usage: deis certs:list [options]

Options:
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

	return cmdr.CertsList(results, time.Now())
}

func certAdd(argv []string, cmdr cmd.Commander) error {
	usage := `
Binds a certificate/key pair to an application.

Usage: deis certs:add <name> <cert> <key> [options]

Arguments:
  <name>
    Name of the certificate to reference it by.
  <cert>
    The public key of the SSL certificate.
  <key>
    The private key of the SSL certificate.

Options:
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	name := args["<name>"].(string)
	cert := args["<cert>"].(string)
	key := args["<key>"].(string)

	return cmdr.CertAdd(cert, key, name)
}

func certRemove(argv []string, cmdr cmd.Commander) error {
	usage := `
removes a certificate/key pair from the application.

Usage: deis certs:remove <name> [options]

Arguments:
  <name>
    the name of the cert to remove from the app.

Options:
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.CertRemove(safeGetValue(args, "<name>"))
}

func certInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
fetch more detailed information about a certificate

Usage: deis certs:info <name> [options]

Arguments:
  <name>
    the name of the cert to get information from

Options:
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.CertInfo(safeGetValue(args, "<name>"))
}

func certAttach(argv []string, cmdr cmd.Commander) error {
	usage := `
attach a certificate to a domain.

Usage: deis certs:attach <name> <domain> [options]

Arguments:
  <name>
    name of the certificate to attach domain to
  <domain>
    common name of the domain to attach to (needs to already be in the system)

Options:
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	name := safeGetValue(args, "<name>")
	domain := safeGetValue(args, "<domain>")
	return cmdr.CertAttach(name, domain)
}

func certDetach(argv []string, cmdr cmd.Commander) error {
	usage := `
detach a certificate from a domain.

Usage: deis certs:detach <name> <domain> [options]

Arguments:
  <name>
    name of the certificate to deatch from a domain
  <domain>
    common name of the domain to detach from

Options:
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	name := safeGetValue(args, "<name>")
	domain := safeGetValue(args, "<domain>")
	return cmdr.CertDetach(name, domain)
}
