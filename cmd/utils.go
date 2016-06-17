package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/workflow-cli/pkg/git"
	"github.com/deis/workflow-cli/settings"
)

var defaultLimit = -1

func progress() chan bool {
	frames := []string{"...", "o..", ".o.", "..o"}
	backspaces := strings.Repeat("\b", 3)
	tick := time.Tick(400 * time.Millisecond)
	quit := make(chan bool)
	go func() {
		for {
			for _, frame := range frames {
				fmt.Print(frame)
				select {
				case <-quit:
					fmt.Print(backspaces)
					close(quit)
					return
				case <-tick:
					fmt.Print(backspaces)
				}
			}
		}
	}()
	return quit
}

// Choose an ANSI color by converting a string to an int.
func chooseColor(input string) string {
	var sum uint8

	for _, char := range []byte(input) {
		sum += uint8(char)
	}

	// Seven possible terminal colors
	color := (sum % 7) + 1

	if color == 7 {
		color = 9
	}

	return fmt.Sprintf("\033[3%dm", color)
}

func load(appID string) (*deis.Client, string, error) {
	c, err := settings.Load()

	if err != nil {
		return nil, "", err
	}

	if appID == "" {
		appID, err = git.DetectAppName(c.ControllerURL.Host)

		if err != nil {
			return nil, "", err
		}
	}

	return c, appID, nil
}

func drinkOfChoice() string {
	drink := os.Getenv("DEIS_DRINK_OF_CHOICE")

	if drink == "" {
		drink = "coffee"
	}

	return drink
}

func limitCount(objs, total int) string {
	if objs == total {
		return "\n"
	}

	return fmt.Sprintf(" (%d of %d)\n", objs, total)
}

// checkAPICompatibility handles specific behavior for certain errors,
// such as printing an warning for the API mismatch error
func checkAPICompatibility(c *deis.Client, err error) error {
	if err == deis.ErrAPIMismatch {
		fmt.Printf(`!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: %s
!    Server version: %s
`, deis.APIVersion, c.ControllerAPIVersion)

		// API mismatch isn't fatal, so after warning continue on.
		return nil
	}

	return err
}
