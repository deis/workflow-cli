package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/deis/workflow-cli/parser"
	docopt "github.com/docopt/docopt-go"
)

// main exits with the return value of Command(os.Args[1:]), deferring all logic to
// a func we can test.
func main() {
	os.Exit(Command(os.Args[1:]))
}

// Command routes deis commands to their proper parser.
func Command(argv []string) int {
	usage := `
The Deis command-line client issues API calls to a Deis controller.

Usage: deis <command> [<args>...]

Option flags::

  -h --help     display help information
  -v --version  display client version

Auth commands::

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Subcommands, use 'deis help [subcommand]' to learn more::

  apps          manage applications used to provide services
  builds        manage builds created using 'git push'
  certs         manage SSL endpoints for an app
  config        manage environment variables that define app config
  domains       manage and assign domain names to your applications
  git           manage git for applications
  keys          manage ssh keys used for 'git push' deployments
  limits        manage resource limits for your application
  perms         manage permissions for applications
  ps            manage processes inside an app container
  registry      manage private registry information for your application
  releases      manage releases of an application
  tags          manage tags for application containers
  users         manage users
  version       display client version

Shortcut commands, use 'deis shortcuts' to see all::

  create        create a new application
  destroy       destroy an application
  info          view information about the current app
  logs          view aggregated log info for the app
  open          open a URL to the app in a browser
  pull          imports an image and deploys as a new release
  run           run a command in an ephemeral app container
  scale         scale processes by type (web=2, worker=1)

Use 'git push deis master' to deploy to an application.
`
	// Reorganize some command line flags and commands.
	command, argv := parseArgs(argv)
	// Give docopt an optional final false arg so it doesn't call os.Exit().
	_, err := docopt.Parse(usage, []string{command}, false, "", true, false)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: deis <command> [<args>...]")
		return 1
	}

	// Dispatch the command, passing the argv through so subcommands can
	// re-parse it according to their usage strings.
	switch command {
	case "apps":
		err = parser.Apps(argv)
	case "auth":
		err = parser.Auth(argv)
	case "builds":
		err = parser.Builds(argv)
	case "certs":
		err = parser.Certs(argv)
	case "config":
		err = parser.Config(argv)
	case "domains":
		err = parser.Domains(argv)
	case "git":
		err = parser.Git(argv)
	case "help":
		fmt.Print(usage)
		return 0
	case "keys":
		err = parser.Keys(argv)
	case "limits":
		err = parser.Limits(argv)
	case "perms":
		err = parser.Perms(argv)
	case "ps":
		err = parser.Ps(argv)
	case "registry":
		err = parser.Registry(argv)
	case "releases":
		err = parser.Releases(argv)
	case "tags":
		err = parser.Tags(argv)
	case "users":
		err = parser.Users(argv)
	case "version":
		err = parser.Version(argv)
	case "doctor":
		err = parser.Doctor(argv)
	default:
		env := os.Environ()
		extCmd := "deis-" + command

		binary, err := exec.LookPath(extCmd)
		if err != nil {
			parser.PrintUsage()
			return 1
		}

		cmdArgv := []string{extCmd}

		cmdSplit := strings.Split(argv[0], command+":")

		if len(cmdSplit) > 1 {
			argv[0] = cmdSplit[1]
		}

		cmdArgv = append(cmdArgv, argv...)

		err = syscall.Exec(binary, cmdArgv, env)
		if err != nil {
			parser.PrintUsage()
			return 1
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

// parseArgs returns the provided args with "--help" as the last arg if need be,
// expands shortcuts and formats commands to be properly routed.
func parseArgs(argv []string) (string, []string) {
	if len(argv) == 1 {
		if argv[0] == "--help" || argv[0] == "-h" {
			// rearrange "deis --help" as "deis help"
			argv[0] = "help"
		} else if argv[0] == "--version" || argv[0] == "-v" {
			// rearrange "deis --version" as "deis version"
			argv[0] = "version"
		}
	}

	if len(argv) >= 2 {
		// Rearrange "deis help <command>" to "deis <command> --help".
		if argv[0] == "help" || argv[0] == "--help" || argv[0] == "-h" {
			argv = append(argv[1:], "--help")
		}
	}

	if len(argv) > 0 {
		argv[0] = replaceShortcut(argv[0])

		index := strings.Index(argv[0], ":")

		if index != -1 {
			command := argv[0]
			return command[:index], argv
		}

		return argv[0], argv
	}

	return "", argv
}

func replaceShortcut(command string) string {
	shortcuts := map[string]string{
		"create":         "apps:create",
		"destroy":        "apps:destroy",
		"info":           "apps:info",
		"login":          "auth:login",
		"logout":         "auth:logout",
		"logs":           "apps:logs",
		"open":           "apps:open",
		"passwd":         "auth:passwd",
		"pull":           "builds:create",
		"register":       "auth:register",
		"rollback":       "releases:rollback",
		"run":            "apps:run",
		"scale":          "ps:scale",
		"sharing":        "perms:list",
		"sharing:list":   "perms:list",
		"sharing:add":    "perms:create",
		"sharing:remove": "perms:delete",
		"whoami":         "auth:whoami",
	}

	expandedCommand := shortcuts[command]
	if expandedCommand == "" {
		return command
	}

	return expandedCommand
}
