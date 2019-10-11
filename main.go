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

	"github.com/bitrise-io/go-steputils/cache"
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
	UseCache    string `env:"cache_local_deps,opt[yes,no]"`
	IsDebugLog  string `env:"verbose_log,opt[yes,no]"`
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

func cacheYarn(workingDir string) error {
	yarnCache := cache.New()
	var cachePaths []string

	// Supporting yarn workspaces (https://yarnpkg.com/lang/en/docs/workspaces/), for this recursively look
	// up all node_modules directories
	if err := filepath.Walk(workingDir, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileInfo.IsDir() && fileInfo.Name() == "node_modules" {
			cachePaths = append(cachePaths, path)
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		failf("Failed to find node_modules directories, error: %s", err)
	}

	log.Debugf("Cached paths: %s", cachePaths)
	for _, path := range cachePaths {
		// yarnCache.IncludePath(fmt.Sprintf("%s -> %s", path, filepath.Join(absWorkingDir, "yarn.lock")))
		yarnCache.IncludePath(path)
	}

	if err := yarnCache.Commit(); err != nil {
		return fmt.Errorf("failed to mark node_modeules directories to be cached, error: %s", err)
	}
	return nil
}

func main() {
	var config config
	if err := stepconf.Parse(&config); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(config)
	fmt.Println()
	log.SetEnableDebugLog(config.IsDebugLog == "yes")

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

	if config.UseCache == "yes" && (len(commandParams) == 0 || commandParams[0] == "install") {
		if err := cacheYarn(absWorkingDir); err != nil {
			failf("%s", err)
		}
	}
}
