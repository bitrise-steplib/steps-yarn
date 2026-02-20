package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-yarn/pkgmanager"
	"github.com/kballard/go-shellquote"
)

type config struct {
	WorkingDir  string `env:"workdir,dir"`
	YarnCommand string `env:"command"`
	YarnArgs    string `env:"args"`
	IsDebugLog  bool   `env:"verbose_log,opt[yes,no]"`
}

func main() {
	var config config
	envRepo := env.NewRepository()
	parser := stepconf.NewInputParser(envRepo)

	logger := log.NewLogger()

	if err := parser.Parse(&config); err != nil {
		failfWithLogger(logger, "Process config: %s", err)
	}
	stepconf.Print(config)

	if config.IsDebugLog {
		logger.EnableDebugLog(true)
	}
	cmdFactory := command.NewFactory(envRepo)
	fmt.Println()

	absWorkingDir, err := filepath.Abs(config.WorkingDir)
	if err != nil {
		failfWithLogger(logger, "Process config: failed to normalize working directory: %s", err)
	}

	commandParams, err := shellquote.Split(config.YarnCommand)
	if err != nil {
		failfWithLogger(logger, "Process config: provided yarn command is not a valid CLI command: %s", err)
	}

	args, err := shellquote.Split(config.YarnArgs)
	if err != nil {
		failfWithLogger(logger, "Process config: provided yarn arguments are not valid CLI arguments: %s", err)
	}

	step := pkgmanager.NewStep(absWorkingDir, logger, cmdFactory)

	if err := step.ValidateAndPrintYarnInfo(absWorkingDir); err != nil {
		failfWithLogger(logger, "Validate Yarn: %s", err)
	}

	if err := step.ExecuteYarnCommand(absWorkingDir, commandParams, args); err != nil {
		failfWithLogger(logger, "Run: %s", err)
	}
}

func failfWithLogger(logger log.Logger, format string, v ...interface{}) {
	logger.Errorf(format, v...)
	os.Exit(1)
}
