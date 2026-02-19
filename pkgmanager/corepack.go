package pkgmanager

import (
	"fmt"
	"strings"
)

// minVersion is the minimum corepack version required to avoid the registry key rotation issue.
// See: https://github.com/nodejs/corepack/issues/612
const minVersion = "0.31.0"

// ensureUpToDate upgrades corepack if the installed version is below minVersion.
// Older versions have stale registry signing keys that cause verification failures
// when downloading package manager versions.
func (s *Step) ensureUpToDate() error {
	cmd := s.cmdFactory.Create("corepack", []string{"--version"}, nil)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err == nil {
		current := strings.TrimSpace(out)
		ok, vErr := versionAtLeast(current, minVersion)
		if vErr != nil {
			s.logger.Debugf("Failed to parse corepack version %q: %s, will upgrade", current, vErr)
		} else if ok {
			s.logger.Infof("Corepack %s is up to date (>= %s)", current, minVersion)
			return nil
		} else {
			s.logger.Infof("Corepack %s is outdated (< %s), upgrading", current, minVersion)
		}
	} else {
		s.logger.Infof("Corepack not found, installing")
	}

	installCmd := s.cmdFactory.Create("npm", []string{"install", "-g", "corepack@latest"}, nil)
	s.logger.Donef("$ %s", installCmd.PrintableCommandArgs())
	if out, err := installCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("failed to install corepack: %s", out)
	}

	return nil
}

// enable runs `corepack enable` to enable all package managers.
func (s *Step) enable() error {
	cmd := s.cmdFactory.Create("corepack", []string{"enable"}, nil)
	s.logger.Donef("$ %s", cmd.PrintableCommandArgs())
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("corepack enable failed: %s", out)
	}

	return nil
}

// versionAtLeast returns true if version >= min (both "X.Y.Z" format).
// Returns false on parse errors, which is a safe fallback (triggers upgrade).
func versionAtLeast(version, min string) (bool, error) {
	var vMaj, vMin, vPat int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &vMaj, &vMin, &vPat); err != nil {
		return false, fmt.Errorf("parse %q: %w", version, err)
	}
	var mMaj, mMin, mPat int
	if _, err := fmt.Sscanf(min, "%d.%d.%d", &mMaj, &mMin, &mPat); err != nil {
		return false, fmt.Errorf("parse %q: %w", min, err)
	}

	if vMaj != mMaj {
		return vMaj > mMaj, nil
	}
	if vMin != mMin {
		return vMin > mMin, nil
	}
	return vPat >= mPat, nil
}
