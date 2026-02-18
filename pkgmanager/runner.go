package pkgmanager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
)

func (s *Step) ExecuteYarnCommand(workDir string, commandParams []string, args []string) error {
	var output bytes.Buffer
	yarnCmd := s.cmdFactory.Create("yarn", append(commandParams, args...), &command.Opts{
		Dir:    workDir,
		Stdout: io.MultiWriter(os.Stdout, &output),
		Stderr: io.MultiWriter(os.Stderr, &output),
	})

	fmt.Println()
	s.logger.Donef("$ %s", yarnCmd.PrintableCommandArgs())
	fmt.Println()

	if err := yarnCmd.Run(); err != nil {
		var exitErr *command.ExitStatusError
		if errors.As(err, &exitErr) {
			if s.isNetworkError(output.String()) {
				fmt.Println()
				s.logger.Warnf(`Looks like you've got network issues while installing yarn.
Please try to increase the timeout with --registry https://registry.npmjs.org --network-timeout [NUMBER] command before using this step (recommended value is 100000).
If issue still persists, please try to debug the error or reach out to support.`)
			}
			return fmt.Errorf("yarn command failed: %s", err)
		}
		return fmt.Errorf("failed to run yarn command: %s", err)
	}

	return nil
}

func (s *Step) isNetworkError(output string) bool {
	return strings.Contains(output, "There appears to be trouble with your network connection. Retrying...")
}
