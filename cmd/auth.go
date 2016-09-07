package cmd

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/auth"
	"github.com/deis/workflow-cli/settings"
	"golang.org/x/crypto/ssh/terminal"
)

// Register creates a account on a Deis controller.
func (d *DeisCmd) Register(controller string, username string, password string, email string,
	sslVerify bool) error {

	c, err := deis.New(sslVerify, controller, "")

	if err != nil {
		return err
	}

	tempSettings, err := settings.Load(d.ConfigFile)

	if err == nil && tempSettings.Client.ControllerURL.Host == c.ControllerURL.Host {
		c.Token = tempSettings.Client.Token
	}

	// Set user agent for temporary client.
	c.UserAgent = settings.UserAgent

	if err = c.CheckConnection(); d.checkAPICompatibility(c, err) != nil {
		return err
	}

	if username == "" {
		d.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		d.Print("password: ")
		password, err = readPassword()
		d.Printf("\npassword (confirm): ")
		passwordConfirm, err := readPassword()
		d.Println()

		if err != nil {
			return err
		}

		if password != passwordConfirm {
			return errors.New("Password mismatch, aborting registration.")
		}
	}

	if email == "" {
		d.Print("email: ")
		fmt.Scanln(&email)
	}

	err = auth.Register(c, username, password, email)

	c.Token = ""

	if d.checkAPICompatibility(c, err) != nil {
		d.PrintErr("Registration failed: ")
		return err
	}

	d.Printf("Registered %s\n", username)

	s := settings.Settings{Client: c}
	return d.doLogin(s, username, password)
}

func (d *DeisCmd) doLogin(s settings.Settings, username, password string) error {
	token, err := auth.Login(s.Client, username, password)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	s.Client.Token = token
	s.Username = username

	filename, err := s.Save(d.ConfigFile)

	if err != nil {
		return nil
	}

	d.Printf("Logged in as %s\n", username)
	d.Printf("Configuration file written to %s\n", filename)
	return nil
}

// Login to a Deis controller.
func (d *DeisCmd) Login(controller string, username string, password string, sslVerify bool) error {
	c, err := deis.New(sslVerify, controller, "")

	if err != nil {
		return err
	}

	// Set user agent for temporary client.
	c.UserAgent = settings.UserAgent

	if err = c.CheckConnection(); d.checkAPICompatibility(c, err) != nil {
		return err
	}

	if username == "" {
		d.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		d.Print("password: ")
		password, err = readPassword()
		d.Println()

		if err != nil {
			return err
		}
	}

	s := settings.Settings{Client: c}
	return d.doLogin(s, username, password)
}

// Logout from a Deis controller.
func (d *DeisCmd) Logout() error {
	if err := settings.Delete(d.ConfigFile); err != nil {
		return err
	}

	d.Println("Logged out")
	return nil
}

// Passwd changes a user's password.
func (d *DeisCmd) Passwd(username, password, newPassword string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if password == "" && username == "" {
		d.Print("current password: ")
		password, err = readPassword()
		d.Println()

		if err != nil {
			return err
		}
	}

	if newPassword == "" {
		d.Print("new password: ")
		newPassword, err = readPassword()
		d.Printf("\nnew password (confirm): ")
		passwordConfirm, err := readPassword()

		d.Println()

		if err != nil {
			return err
		}

		if newPassword != passwordConfirm {
			return errors.New("Password mismatch, not changing.")
		}
	}

	err = auth.Passwd(s.Client, username, password, newPassword)
	if d.checkAPICompatibility(s.Client, err) != nil {
		d.PrintErr("Password change failed: ")
		return err
	}

	d.Println("Password change succeeded.")
	return nil
}

// Cancel deletes a user's account.
func (d *DeisCmd) Cancel(username, password string, yes bool) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if username == "" || password != "" {
		d.Println("Please log in again in order to cancel this account")

		if err = d.Login(s.Client.ControllerURL.String(), username, password, s.Client.VerifySSL); err != nil {
			return err
		}
	}

	if !yes {
		confirm := ""

		s, err = settings.Load(d.ConfigFile)

		if err != nil {
			return err
		}

		deletedUser := username

		if deletedUser == "" {
			deletedUser = s.Username
		}

		d.Printf("cancel account %s at %s? (y/N): ", deletedUser, s.Client.ControllerURL.String())
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "y" {
			yes = true
		}
	}

	if !yes {
		d.PrintErrln("Account not changed")
		return nil
	}

	err = auth.Delete(s.Client, username)
	if err == deis.ErrConflict {
		return fmt.Errorf("%s still has applications associated with it. Transfer ownership or delete them first", username)
	} else if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	// If user targets themselves, logout.
	if username == "" || s.Username == username {
		if err := settings.Delete(d.ConfigFile); err != nil {
			return err
		}
	}

	d.Println("Account cancelled")
	return nil
}

// Whoami prints the logged in user. If all is true, it fetches info from the controller to know
// more about the user.
func (d *DeisCmd) Whoami(all bool) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if all {
		user, err := auth.Whoami(s.Client)
		if err != nil {
			return err
		}
		d.Println(user)
	} else {
		d.Printf("You are %s at %s\n", s.Username, s.Client.ControllerURL.String())
	}
	return nil
}

// Regenerate regenenerates a user's token.
func (d *DeisCmd) Regenerate(username string, all bool) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	token, err := auth.Regenerate(s.Client, username, all)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if username == "" && !all {
		s.Client.Token = token
		_, err = s.Save(d.ConfigFile)

		if err != nil {
			return err
		}
	}

	d.Println("Token Regenerated")
	return nil
}

func readPassword() (string, error) {
	password, err := terminal.ReadPassword(int(syscall.Stdin))

	return string(password), err
}
