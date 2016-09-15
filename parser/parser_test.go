package parser

import "io"

// DeisCmd is an implementation of Commander.
type FakeDeisCmd struct {
	ConfigFile string
	Warned     bool
	WOut       io.Writer
	WErr       io.Writer
	WIn        io.Reader
}

func (d FakeDeisCmd) Println(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDeisCmd) Print(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDeisCmd) Printf(string, ...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDeisCmd) PrintErrln(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDeisCmd) PrintErr(...interface{}) (int, error) {
	return 1, nil
}

func (d FakeDeisCmd) PrintErrf(string, ...interface{}) (int, error) {
	return 1, nil
}
