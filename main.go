package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/kballard/go-shellquote"
)

type config struct {
	WorkingDir  string `env:"workdir,dir"`
	YarnCommand string `env:"command"`
	YarnArgs    string `env:"args"`
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func getInstallYarnCommand() (*command.Model, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("Unsupported platform %s, failed to install yarn", runtime.GOOS)
	}
	if _, err := os.Stat(path.Join("etc", "lsb-release")); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("only Ubuntu distribution supported, failed to install yarn")
		}
		return nil, err
	}

	installCmd := command.New("sh", "-c", `curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt-get update && sudo apt-get install -y yarn`)
	installCmd.SetStdout(os.Stdout).SetStderr(os.Stderr)

	return installCmd, nil
}

func main() {
	var config config
	if err := stepconf.Parse(&config); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(config)
	fmt.Println()

	if path, err := exec.LookPath("yarn"); err != nil {
		log.Infof("Yarn not installed. Installing...")
		installCmd, err := getInstallYarnCommand()
		if err != nil {
			failf("%s", err)
		}

		fmt.Println()
		log.Donef("$ %s", installCmd.PrintableCommandArgs())
		fmt.Println()

		if err := installCmd.Run(); err != nil {
			if errorutil.IsExitStatusError(err) {
				failf("yarn install command failed, error: %s", err)
			}
			failf("Failed to run command, error: %s", err)
		}
	} else {
		log.Infof("Yarn is already installed at: %s", path)
	}

	absWorkingDir, err := filepath.Abs(config.WorkingDir)
	if err != nil {
		failf("Failed to get absolute working directory, error: %s", err)
	}

	log.Infof("Yarn version:")
	versionCmd := command.New("yarn", "--version")
	versionCmd.SetStdout(os.Stdout).SetStderr(os.Stderr).SetDir(absWorkingDir)

	fmt.Println()
	log.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()
	err = versionCmd.Run()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			failf("yarn version command failed, error: %s", err)
		}
		failf("Failed to run command, error: %s", err)
	}

	commandParams, err := shellquote.Split(config.YarnCommand)
	if err != nil {
		failf("failed to split command arguments, error: %s", err)
	}

	args, err := shellquote.Split(config.YarnArgs)
	if err != nil {
		failf("failed to split command arguments, error: %s", err)
	}

	yarnCmd := command.New("yarn", append(commandParams, args...)...)
	var output bytes.Buffer
	outputWriter := io.MultiWriter(os.Stdout, &output)
	yarnCmd.SetDir(absWorkingDir).SetStdout(outputWriter).SetStderr(outputWriter)

	fmt.Println()
	log.Donef("$ %s", yarnCmd.PrintableCommandArgs())
	fmt.Println()

	defer func() {
		if strings.Contains(output.String(), "There appears to be trouble with your network connection. Retrying...") {
			fmt.Println()
			log.Warnf(`Looks like you've got network issues while installing yarn.
Please try to increase the timeout with --registry https://registry.npmjs.org --network-timeout [NUMBER] command before using this step (recommended value is 100000).
If issue still persists, please try to debug the error or reach out to support.`)
		}
	}()

	if err := yarnCmd.Run(); err != nil {
		if errorutil.IsExitStatusError(err) {
			failf("yarn command failed, error: %s", err)
		}
		failf("failed to run yarn command, error: %s", err)
	}
}
