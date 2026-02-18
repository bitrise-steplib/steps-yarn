package pkgmanager

import (
	"fmt"

	"github.com/bitrise-io/go-utils/v2/command"
)

func (s *Step) ensureCorepack(workDir string) error {
	s.logger.Infof("Checking Corepack status...")
	versionCmd := s.cmdFactory.Create("corepack", []string{"--version"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	s.logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()

	if out, err := versionCmd.RunAndReturnTrimmedCombinedOutput(); err == nil {
		s.logger.Printf("Corepack version: %s", out)
		fmt.Println()
		return nil
	}

	s.logger.Infof("Could not verify Corepack, installing...")
	installCmd := s.cmdFactory.Create("npm", []string{"install", "--global", "corepack"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	s.logger.Donef("$ %s", installCmd.PrintableCommandArgs())
	fmt.Println()

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("install Corepack: %s", err)
	}

	fmt.Println()
	s.logger.Donef("$ %s", versionCmd.PrintableCommandArgs())
	fmt.Println()

	if out, err := versionCmd.RunAndReturnTrimmedCombinedOutput(); err == nil {
		s.logger.Infof("Corepack installed successfully")
		s.logger.Printf("Corepack version: %s", out)
	} else {
		return fmt.Errorf("verify Corepack installation: %s", err)
	}
	fmt.Println()

	s.logger.Infof("Enable Corepack...")
	enableCmd := s.cmdFactory.Create("corepack", []string{"enable"}, &command.Opts{
		Dir: workDir,
	})

	fmt.Println()
	s.logger.Donef("$ %s", enableCmd.PrintableCommandArgs())
	fmt.Println()

	if err := enableCmd.Run(); err != nil {
		return fmt.Errorf("enable Corepack: %s", err)
	}

	s.logger.Infof("Corepack enabled successfully")
	fmt.Println()

	return nil
}
