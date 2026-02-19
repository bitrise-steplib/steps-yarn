package pkgmanager

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

// minVersion is the minimum corepack version required to avoid the registry key rotation issue.
// See: https://github.com/nodejs/corepack/issues/612
const minVersion = "0.31.0"

// EnsureUpToDate upgrades corepack if the installed version is below minVersion.
// Older versions have stale registry signing keys that cause verification failures
// when downloading package manager versions.
func EnsureUpToDate(cmdFactory command.Factory, logger log.Logger) error {
	if out, err := exec.Command("corepack", "--version").Output(); err == nil {
		current := strings.TrimSpace(string(out))
		if ok, _ := versionAtLeast(current, minVersion); ok {
			logger.Infof("Corepack %s is up to date (>= %s)", current, minVersion)
			return nil
		}
		logger.Infof("Corepack %s is outdated (< %s), upgrading", current, minVersion)
	} else {
		logger.Infof("Corepack not found, installing")
	}

	cmd := cmdFactory.Create("npm", []string{"install", "-g", "corepack@latest"}, nil)
	logger.Donef("$ %s", cmd.PrintableCommandArgs())
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("failed to install corepack: %s", out)
	}

	return nil
}

// Enable runs `corepack enable` to enable all package managers.
func Enable(cmdFactory command.Factory, logger log.Logger) error {
	cmd := cmdFactory.Create("corepack", []string{"enable"}, nil)
	logger.Donef("$ %s", cmd.PrintableCommandArgs())
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
