package pkgmanager

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log"
)

func TestEnsureUpToDate(t *testing.T) {
	tests := []struct {
		name          string
		runFunc       func(cmdStr string) (string, error)
		expectError   bool
		errorContains string
	}{
		{
			name: "corepack up to date",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "0.31.0", nil
				}
				return "", nil
			},
			expectError: false,
		},
		{
			name: "corepack outdated, upgrade succeeds",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "0.30.0", nil
				}
				if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
					return "added corepack@latest", nil
				}
				return "", nil
			},
			expectError: false,
		},
		{
			name: "corepack not installed, install succeeds",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "", &mockExitError{}
				}
				if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
					return "added corepack@latest", nil
				}
				return "", nil
			},
			expectError: false,
		},
		{
			name: "npm install fails",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "", &mockExitError{}
				}
				if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
					return "network error", &mockExitError{}
				}
				return "", nil
			},
			expectError:   true,
			errorContains: "failed to install corepack",
		},
		{
			name: "invalid version format triggers upgrade",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "invalid.version", nil
				}
				if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
					return "added corepack@latest", nil
				}
				return "", nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.NewLogger(log.WithOutput(&logBuf))

			factory := &mockFactory{runFunc: tt.runFunc}
			step := &Step{
				logger:     logger,
				cmdFactory: factory,
			}

			err := step.ensureUpToDate()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.errorContains)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.errorContains, err)
				}
			} else if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestEnable(t *testing.T) {
	tests := []struct {
		name          string
		runFunc       func(cmdStr string) (string, error)
		expectError   bool
		errorContains string
	}{
		{
			name: "enable succeeds",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "enable") {
					return "", nil
				}
				return "", nil
			},
			expectError: false,
		},
		{
			name: "enable fails",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "enable") {
					return "permission denied", &mockExitError{}
				}
				return "", nil
			},
			expectError:   true,
			errorContains: "corepack enable failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.NewLogger(log.WithOutput(&logBuf))

			factory := &mockFactory{runFunc: tt.runFunc}
			step := &Step{
				logger:     logger,
				cmdFactory: factory,
			}

			err := step.enable()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.errorContains)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.errorContains, err)
				}
			} else if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestVersionAtLeast(t *testing.T) {
	tests := []struct {
		name    string
		version string
		min     string
		want    bool
	}{
		{"equal versions", "0.31.0", "0.31.0", true},
		{"newer patch", "0.31.1", "0.31.0", true},
		{"older patch", "0.31.0", "0.31.1", false},
		{"newer minor", "0.32.0", "0.31.0", true},
		{"older minor", "0.30.0", "0.31.0", false},
		{"newer major", "1.0.0", "0.31.0", true},
		{"older major", "0.31.0", "1.0.0", false},
		{"complex newer", "1.2.3", "0.31.0", true},
		{"complex older", "0.30.9", "0.31.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := versionAtLeast(tt.version, tt.min)
			if err != nil {
				t.Errorf("versionAtLeast() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("versionAtLeast(%q, %q) = %v, want %v", tt.version, tt.min, got, tt.want)
			}
		})
	}
}

func TestVersionAtLeast_InvalidVersions(t *testing.T) {
	tests := []struct {
		name    string
		version string
		min     string
	}{
		{"invalid version", "invalid", "0.31.0"},
		{"invalid min", "0.31.0", "invalid"},
		{"incomplete version", "0.31", "0.31.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := versionAtLeast(tt.version, tt.min)
			if err == nil {
				t.Errorf("versionAtLeast(%q, %q) expected error, got result: %v", tt.version, tt.min, got)
			}
			if got {
				t.Errorf("versionAtLeast(%q, %q) should return false on error, got true", tt.version, tt.min)
			}
		})
	}
}

type mockExitError struct{}

func (e *mockExitError) Error() string {
	return "exit status 1"
}
