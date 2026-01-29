package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/kballard/go-shellquote"
)

type config struct {
	WorkingDir  string `env:"workdir,dir"`
	YarnCommand string `env:"command"`
	YarnArgs    string `env:"args"`
	IsDebugLog  bool   `env:"verbose_log,opt[yes,no]"`
}

var (
	logger     log.Logger
	cmdFactory command.Factory
)

func main() {
	var config config
	envRepo := env.NewRepository()
	parser := stepconf.NewInputParser(envRepo)
	if err := parser.Parse(&config); err != nil {
		failf("Process config: %s", err)
	}
	stepconf.Print(config)

	logger = log.NewLogger()
	if config.IsDebugLog {
		logger.EnableDebugLog(true)
	}
	cmdFactory = command.NewFactory(envRepo)
	fmt.Println()

	absWorkingDir, err := filepath.Abs(config.WorkingDir)
	if err != nil {
		failf("Process config: failed to normalize working directory: %s", err)
	}

	commandParams, err := shellquote.Split(config.YarnCommand)
	if err != nil {
		failf("Process config: provided yarn command is not a valid CLI command: %s", err)
	}

	args, err := shellquote.Split(config.YarnArgs)
	if err != nil {
		failf("Process config: provided yarn arguments are not valid CLI arguments: %s", err)
	}

	validInstallation := validateYarnInstallation(absWorkingDir)
	if !validInstallation {
		if err := installYarn(); err != nil {
			failf("Install dependencies: %s", err)
		}
		if err := printYarnVersion(absWorkingDir); err != nil {
			failf("Install dependencies: %s", err)
		}
	}

	var output bytes.Buffer
	yarnCmd := cmdFactory.Create("yarn", append(commandParams, args...), &command.Opts{
		Dir:    absWorkingDir,
		Stdout: io.MultiWriter(os.Stdout, &output),
		Stderr: io.MultiWriter(os.Stderr, &output),
	})

	fmt.Println()
	logger.Donef("$ %s", yarnCmd.PrintableCommandArgs())
	fmt.Println()

	if err := yarnCmd.Run(); err != nil {
		var exitErr *command.ExitStatusError
		if errors.As(err, &exitErr) {
			if strings.Contains(output.String(), "There appears to be trouble with your network connection. Retrying...") {
				fmt.Println()
				logger.Warnf(`Looks like you've got network issues while installing yarn.
	Please try to increase the timeout with --registry https://registry.npmjs.org --network-timeout [NUMBER] command before using this step (recommended value is 100000).
	If issue still persists, please try to debug the error or reach out to support.`)
			}
			failf("Run: provided yarn command failed: %s", err)
		}
		failf("Run: failed to run provided yarn command: %s", err)
	}
}

func failf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
	os.Exit(1)
}

func getInstallYarnCommand() command.Command {
	return cmdFactory.Create("npm", []string{"install", "--global", "yarn"}, nil)
}


func validateYarnInstallation(workDir string) bool {
	pth, err := exec.LookPath("yarn")
	if err != nil {
		logger.Debugf("yarn is not installed to the PATH")
		return false
	}

	versionCmd := cmdFactory.Create("yarn", []string{"--version"}, &command.Opts{
		Dir: workDir,
	})
	out, err := versionCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		logger.Debugf("yarn version command failed: %s, out: %s", err, out)
		return false
	}

	logger.Infof("Yarn is already installed at: %s", pth)
	fmt.Println()
	logger.Infof("Yarn version:")
	logger.Printf(out)

	return true
}

func installYarn() error {
	logger.Infof("Yarn not installed. Installing...")
	installCmd := getInstallYarnCommand()

	fmt.Println()
	logger.Donef("$ %s", installCmd.PrintableCommandArgs())
	fmt.Println()

	if err := installCmd.Run(); err != nil {
		var exitErr *command.ExitStatusError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("installing yarn failed: %s", err)
		}
		return fmt.Errorf("failed to run command: %s", err)
	}

	return nil
}

func printYarnVersion(workDir string) error {
	logger.Infof("Yarn version:")
	versionCmd := cmdFactory.Create("yarn", []string{"--version"}, &command.Opts{
		Dir:    workDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	fmt.Println()
	logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()
	if err := versionCmd.Run(); err != nil {
		var exitErr *command.ExitStatusError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("yarn version command failed: %s", err)
		}
		return fmt.Errorf("failed to run command: %s", err)
	}

	return nil
}
