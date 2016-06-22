package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/auth"
	"github.com/deis/workflow-cli/settings"
	"golang.org/x/crypto/ssh/terminal"
)

// Register creates a account on a Deis controller.
func Register(controller string, username string, password string, email string,
	sslVerify bool) error {

	c, err := deis.New(sslVerify, controller, "", "")

	if err != nil {
		return err
	}

	tempClient, err := settings.Load()

	if err == nil && tempClient.ControllerURL.Host == c.ControllerURL.Host {
		c.Token = tempClient.Token
	}

	// Set user agent for temporary client.
	c.UserAgent = settings.UserAgent

	if err = c.CheckConnection(); checkAPICompatibility(c, err) != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Printf("\npassword (confirm): ")
		passwordConfirm, err := readPassword()
		fmt.Println()

		if err != nil {
			return err
		}

		if password != passwordConfirm {
			return errors.New("Password mismatch, aborting registration.")
		}
	}

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	err = auth.Register(c, username, password, email)

	c.Token = ""

	if checkAPICompatibility(c, err) != nil {
		fmt.Fprint(os.Stderr, "Registration failed: ")
		return err
	}

	fmt.Printf("Registered %s\n", username)
	return doLogin(c, username, password)
}

func doLogin(c *deis.Client, username, password string) error {
	token, err := auth.Login(c, username, password)
	if checkAPICompatibility(c, err) != nil {
		return err
	}

	c.Token = token
	c.Username = username

	err = settings.Save(c)

	if err != nil {
		return nil
	}

	fmt.Printf("Logged in as %s\n", username)
	return nil
}

// Login to a Deis controller.
func Login(controller string, username string, password string, sslVerify bool) error {
	c, err := deis.New(sslVerify, controller, "", "")

	if err != nil {
		return err
	}

	// Set user agent for temporary client.
	c.UserAgent = settings.UserAgent

	if err = c.CheckConnection(); checkAPICompatibility(c, err) != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	return doLogin(c, username, password)
}

// Logout from a Deis controller.
func Logout() error {
	if err := settings.Delete(); err != nil {
		return err
	}

	fmt.Println("Logged out")
	return nil
}

// Passwd changes a user's password.
func Passwd(username string, password string, newPassword string) error {
	c, err := settings.Load()

	if err != nil {
		return err
	}

	if password == "" && username == "" {
		fmt.Print("current password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	if newPassword == "" {
		fmt.Print("new password: ")
		newPassword, err = readPassword()
		fmt.Printf("\nnew password (confirm): ")
		passwordConfirm, err := readPassword()

		fmt.Println()

		if err != nil {
			return err
		}

		if newPassword != passwordConfirm {
			return errors.New("Password mismatch, not changing.")
		}
	}

	err = auth.Passwd(c, username, password, newPassword)
	if checkAPICompatibility(c, err) != nil {
		fmt.Fprint(os.Stderr, "Password change failed: ")
		return err
	}

	fmt.Println("Password change succeeded.")
	return nil
}

// Cancel deletes a user's account.
func Cancel(username string, password string, yes bool) error {
	c, err := settings.Load()

	if err != nil {
		return err
	}

	if username == "" || password != "" {
		fmt.Println("Please log in again in order to cancel this account")

		if err = Login(c.ControllerURL.String(), username, password, c.VerifySSL); err != nil {
			return err
		}
	}

	if yes == false {
		confirm := ""

		c, err = settings.Load()

		if err != nil {
			return err
		}

		deletedUser := username

		if deletedUser == "" {
			deletedUser = c.Username
		}

		fmt.Printf("cancel account %s at %s? (y/N): ", deletedUser, c.ControllerURL.String())
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "y" {
			yes = true
		}
	}

	if yes == false {
		fmt.Fprintln(os.Stderr, "Account not changed")
		return nil
	}

	err = auth.Delete(c, username)
	if err == deis.ErrConflict {
		return fmt.Errorf("%s still has applications associated with it. Transfer ownership or delete them first", username)
	} else if checkAPICompatibility(c, err) != nil {
		return err
	}

	// If user targets themselves, logout.
	if username == "" || c.Username == username {
		if err := settings.Delete(); err != nil {
			return err
		}
	}

	fmt.Println("Account cancelled")
	return nil
}

// Whoami prints the logged in user.
func Whoami() error {
	c, err := settings.Load()

	if err != nil {
		return err
	}

	fmt.Printf("You are %s at %s\n", c.Username, c.ControllerURL.String())
	return nil
}

// Regenerate regenenerates a user's token.
func Regenerate(username string, all bool) error {
	c, err := settings.Load()

	if err != nil {
		return err
	}

	token, err := auth.Regenerate(c, username, all)
	if checkAPICompatibility(c, err) != nil {
		return err
	}

	if username == "" && all == false {
		c.Token = token

		err = settings.Save(c)

		if err != nil {
			return err
		}
	}

	fmt.Println("Token Regenerated")
	return nil
}

func readPassword() (string, error) {
	password, err := terminal.ReadPassword(int(syscall.Stdin))

	return string(password), err
}
