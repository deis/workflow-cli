package settings

import (
	"os"
	"path/filepath"
)

func locateSettingsFile() string {
	filename := os.Getenv("DEIS_PROFILE")

	if filename == "" {
		filename = "client"
	}

	return filepath.Join(FindHome(), ".deis", filename+".json")
}
