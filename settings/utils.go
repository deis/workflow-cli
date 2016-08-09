package settings

import (
	"os"
	"path/filepath"
	"regexp"
)

var filepathRegex = regexp.MustCompile(`^.*[/\\].+\.json$`)

func locateSettingsFile(cf string) string {
	if cf == "" {
		if v, ok := os.LookupEnv("DEIS_PROFILE"); ok {
			cf = v
		} else {
			cf = "client"
		}
	}

	// if path appears to be a filepath (contains a separator and ends in .json) don't alter the path
	if filepathRegex.MatchString(cf) {
		return cf
	}

	return filepath.Join(FindHome(), ".deis", cf+".json")
}
