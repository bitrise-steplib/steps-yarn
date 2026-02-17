package main

import (
	"bytes"
	"encoding/json"
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

	if err := validateAndPrintYarnInfo(absWorkingDir); err != nil {
		failf("Validate Yarn: %s", err)
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

func validateAndPrintYarnInfo(workDir string) error {
	if err := ensureCorepack(workDir); err != nil {
		return err
	}

	pth, err := exec.LookPath("yarn")
	if err != nil {
		return fmt.Errorf("yarn is unexpectedly not installed to the PATH after Corepack setup: %s", err)
	}

	logger.Infof("Using Yarn from PATH (managed by Corepack): %s", pth)
	fmt.Println()

	printPackageManagerInfo(workDir)

	printYarnrcInfo(workDir)

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
		return fmt.Errorf("yarn version command: %s", err)
	}

	return nil
}

func printPackageManagerInfo(workDir string) {
	packageJSONPath := filepath.Join(workDir, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		logger.Debugf("Could not read package.json: %s", err)
		return
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		logger.Debugf("Could not parse package.json: %s", err)
		return
	}

	if packageManager, ok := packageJSON["packageManager"].(string); ok {
		logger.Infof("package.json specifies packageManager: %s", packageManager)
		fmt.Println()
	} else {
		logger.Debugf("No packageManager field found in package.json")
	}
}

func printYarnrcInfo(workDir string) {
	yarnrcPath := filepath.Join(workDir, ".yarnrc")
	if data, err := os.ReadFile(yarnrcPath); err == nil {
		logger.Infof(".yarnrc file found (Yarn 1.x config):")
		logger.Printf(strings.TrimSpace(string(data)))
		fmt.Println()
		return
	}

	yarnrcYmlPath := filepath.Join(workDir, ".yarnrc.yml")
	if _, err := os.Stat(yarnrcYmlPath); err == nil {
		logger.Infof(".yarnrc.yml file found (Yarn 2+ config)")
		fmt.Println()
	}
}

func ensureCorepack(workDir string) error {
	logger.Infof("Checking Corepack status...")
	versionCmd := cmdFactory.Create("corepack", []string{"--version"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()

	if out, err := versionCmd.RunAndReturnTrimmedCombinedOutput(); err == nil {
		logger.Printf("Corepack version: %s", out)
		fmt.Println()
		return nil
	}

	logger.Infof("Could not verify Corepack, installing...")
	installCmd := cmdFactory.Create("npm", []string{"install", "--global", "corepack"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	logger.Donef("$ %s", installCmd.PrintableCommandArgs())
	fmt.Println()

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("install Corepack: %s", err)
	}

	fmt.Println()
	logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()

	if out, err := versionCmd.RunAndReturnTrimmedCombinedOutput(); err == nil {
		logger.Infof("Corepack installed successfully")
		logger.Printf("Corepack version: %s", out)
	} else {
		return fmt.Errorf("verify Corepack installation: %s", err)
	}
	fmt.Println()

	logger.Infof("Enable Corepack...")
	enableCmd := cmdFactory.Create("corepack", []string{"enable"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	logger.Donef("$ %s", enableCmd.PrintableCommandArgs())
	fmt.Println()

	if err := enableCmd.Run(); err != nil {
		return fmt.Errorf("enable Corepack: %s", err)
	}

	logger.Infof("Corepack enabled successfully")
	fmt.Println()

	return nil
}
