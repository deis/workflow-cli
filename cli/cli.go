package cli

// Shortcuts is a map of all the shortcuts supported by the CLI
var Shortcuts = map[string]string{
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
