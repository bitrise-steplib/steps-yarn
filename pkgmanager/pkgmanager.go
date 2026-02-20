package pkgmanager

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Step struct {
	logger     log.Logger
	cmdFactory command.Factory
	fs         fs.FS
}

func NewStep(workDir string, logger log.Logger, cmdFactory command.Factory) *Step {
	return &Step{
		logger:     logger,
		cmdFactory: cmdFactory,
		fs:         os.DirFS(workDir),
	}
}

func (s *Step) ValidateAndPrintYarnInfo(workDir string) error {
	// Ensure corepack is up to date
	if err := s.ensureUpToDate(); err != nil {
		return err
	}

	if err := s.enable(); err != nil {
		return err
	}

	pth, err := exec.LookPath("yarn")
	if err != nil {
		return fmt.Errorf("yarn is unexpectedly not installed to the PATH after Corepack setup: %s", err)
	}

	s.logger.Infof("Using Yarn from PATH (managed by Corepack): %s", pth)
	fmt.Println()

	s.printPackageManagerInfo()

	s.printYarnrcInfo()

	s.logger.Infof("Yarn version:")
	versionCmd := s.cmdFactory.Create("yarn", []string{"--version"}, &command.Opts{
		Dir:    workDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	fmt.Println()
	s.logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()

	if err := versionCmd.Run(); err != nil {
		return fmt.Errorf("yarn version command: %s", err)
	}

	return nil
}

func (s *Step) printPackageManagerInfo() {
	data, err := fs.ReadFile(s.fs, "package.json")
	if err != nil {
		s.logger.Debugf("Could not read package.json: %s", err)
		return
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		s.logger.Debugf("Could not parse package.json: %s", err)
		return
	}

	if packageManager, ok := packageJSON["packageManager"].(string); ok {
		s.logger.Infof("package.json specifies packageManager: %s", packageManager)
		fmt.Println()
	} else {
		s.logger.Warnf("No packageManager field found in package.json")
		fmt.Println()
	}
}

func (s *Step) printYarnrcInfo() {
	if data, err := fs.ReadFile(s.fs, ".yarnrc"); err == nil {
		s.logger.Infof(".yarnrc file found (Yarn 1.x config):")
		s.logger.Printf(strings.TrimSpace(string(data)))
		fmt.Println()
		return
	}

	if _, err := fs.Stat(s.fs, ".yarnrc.yml"); err == nil {
		s.logger.Infof(".yarnrc.yml file found (Yarn 2+ config)")
		fmt.Println()
	}
}
