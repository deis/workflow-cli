package cmd

import (
	"fmt"
	"io"

	"github.com/deis/controller-sdk-go/api"
)

// Commander is interface definition for running commands
type Commander interface {
	AppCreate(string, string, string, bool) error
	AppsList(int) error
	AppInfo(string) error
	AppOpen(string) error
	AppLogs(string, int) error
	AppRun(string, string) error
	AppDestroy(string, string) error
	AppTransfer(string, string) error
	AutoscaleList(string) error
	AutoscaleSet(string, string, int, int, int) error
	AutoscaleUnset(string, string) error
	Register(string, string, string, string, bool) error
	Login(string, string, string, bool) error
	Logout() error
	Passwd(string, string, string) error
	Cancel(string, string, bool) error
	Whoami(bool) error
	Regenerate(string, bool) error
	BuildsList(string, int) error
	BuildsCreate(string, string, string) error
	CertsList(int) error
	CertAdd(string, string, string) error
	CertRemove(string) error
	CertInfo(string) error
	CertAttach(string, string) error
	CertDetach(string, string) error
	ConfigList(string, bool) error
	ConfigSet(string, []string) error
	ConfigUnset(string, []string) error
	ConfigPull(string, bool, bool) error
	ConfigPush(string, string) error
	DomainsList(string, int) error
	DomainsAdd(string, string) error
	DomainsRemove(string, string) error
	GitRemote(string, string, bool) error
	GitRemove(string) error
	HealthchecksList(string, string) error
	HealthchecksSet(string, string, string, *api.Healthcheck) error
	HealthchecksUnset(string, string, []string) error
	KeysList(int) error
	KeyRemove(string) error
	KeyAdd(string) error
	LimitsList(string) error
	LimitsSet(string, []string, string) error
	LimitsUnset(string, []string, string) error
	MaintenanceInfo(string) error
	MaintenanceEnable(string) error
	MaintenanceDisable(string) error
	PermsList(string, bool, int) error
	PermCreate(string, string, bool) error
	PermDelete(string, string, bool) error
	PsList(string, int) error
	PsScale(string, []string) error
	PsRestart(string, string) error
	RegistryList(string) error
	RegistrySet(string, []string) error
	RegistryUnset(string, []string) error
	ReleasesList(string, int) error
	ReleasesInfo(string, int) error
	ReleasesRollback(string, int) error
	RoutingInfo(string) error
	RoutingEnable(string) error
	RoutingDisable(string) error
	ShortcutsList() error
	TagsList(string) error
	TagsSet(string, []string) error
	TagsUnset(string, []string) error
	UsersList(results int) error
	WhitelistAdd(string, string) error
	WhitelistList(string) error
	WhitelistRemove(string, string) error
	Println(...interface{}) (int, error)
	Print(...interface{}) (int, error)
	Printf(string, ...interface{}) (int, error)
	PrintErrln(...interface{}) (int, error)
	PrintErr(...interface{}) (int, error)
	PrintErrf(string, ...interface{}) (int, error)
}

// DeisCmd is an implementation of Commander.
type DeisCmd struct {
	ConfigFile string
	WOut       io.Writer
	WErr       io.Writer
}

// Println prints a line to an output writer.
func (d *DeisCmd) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(d.WOut, a...)
}

// Print prints a line to an output writer.
func (d *DeisCmd) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(d.WOut, a...)
}

// Printf prints a line to an error writer.
func (d *DeisCmd) Printf(s string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(d.WOut, s, a...)
}

// PrintErrln prints a line to an error writer.
func (d *DeisCmd) PrintErrln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(d.WErr, a...)
}

// PrintErr prints a line to an error writer.
func (d *DeisCmd) PrintErr(a ...interface{}) (n int, err error) {
	return fmt.Fprint(d.WErr, a...)
}

// PrintErrf prints a line to an error writer.
func (d *DeisCmd) PrintErrf(s string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(d.WErr, s, a...)
}
