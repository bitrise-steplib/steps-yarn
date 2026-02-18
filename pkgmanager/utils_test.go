package pkgmanager

import (
	"fmt"

	"github.com/bitrise-io/go-utils/v2/command"
)

type mockCommand struct {
	cmdStr  string
	opts    *command.Opts
	runFunc func(cmdStr string) (output string, err error)
}

func (m *mockCommand) PrintableCommandArgs() string { return m.cmdStr }

func (m *mockCommand) Run() error {
	output, err := m.runFunc(m.cmdStr)
	if output != "" && m.opts != nil && m.opts.Stderr != nil {
		fmt.Fprint(m.opts.Stderr, output)
	}
	return err
}

func (m *mockCommand) RunAndReturnTrimmedCombinedOutput() (string, error) {
	return m.runFunc(m.cmdStr)
}

func (m *mockCommand) RunAndReturnTrimmedOutput() (string, error) {
	return m.runFunc(m.cmdStr)
}

func (m *mockCommand) RunAndReturnExitCode() (int, error) {
	_, err := m.runFunc(m.cmdStr)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func (m *mockCommand) Start() error { return nil }
func (m *mockCommand) Wait() error  { return nil }

type mockFactory struct {
	runFunc func(cmdStr string) (output string, err error)
}

func (m *mockFactory) Create(name string, args []string, opts *command.Opts) command.Command {
	return &mockCommand{
		cmdStr:  fmt.Sprintf("%s %v", name, args),
		opts:    opts,
		runFunc: m.runFunc,
	}
}
