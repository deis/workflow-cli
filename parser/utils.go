package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/deis/workflow-cli/cmd"
)

func safeGetValue(args map[string]interface{}, key string) string {
	if args[key] == nil {
		return ""
	}
	return args[key].(string)
}

func safeGetInt(args map[string]interface{}, key string) int {
	if args[key] == nil {
		return 0
	}
	retVal, err := strconv.Atoi(args[key].(string))
	if err != nil {
		log.Fatalf("could not convert %s to int: %v", args[key], err)
	}
	return retVal
}

func responseLimit(limit string) (int, error) {
	if limit == "" {
		return -1, nil
	}

	return strconv.Atoi(limit)
}

// PrintUsage runs if no matching command is found.
func PrintUsage(cmdr cmd.Commander) {
	cmdr.PrintErrln("Found no matching command, try 'deis help'")
	cmdr.PrintErrln("Usage: deis <command> [<args>...]")
}

func printHelp(argv []string, usage string) bool {
	if len(argv) > 1 {
		if argv[1] == "--help" || argv[1] == "-h" {
			fmt.Print(usage)
			return true
		}
	}

	return false
}
