package pkgmanager

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

func TestEnsureCorepack(t *testing.T) {
	tests := []struct {
		name          string
		runFunc       func(cmdStr string) (string, error)
		expectError   bool
		errorContains string
		expectLogMsg  string
	}{
		{
			name: "corepack already installed",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "0.17.0", nil
				}
				return "", nil
			},
			expectError: false,
		},
		{
			name: "corepack needs installation",
			runFunc: func() func(cmdStr string) (string, error) {
				callCount := 0
				return func(cmdStr string) (string, error) {
					if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
						callCount++
						if callCount == 1 {
							return "", fmt.Errorf("command not found")
						}
						return "0.17.0", nil
					}
					if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
						return "added corepack", nil
					}
					if strings.Contains(cmdStr, "enable") {
						return "", nil
					}
					return "", nil
				}
			}(),
			expectError: false,
		},
		{
			name: "npm install fails",
			runFunc: func(cmdStr string) (string, error) {
				if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
					return "", fmt.Errorf("command not found")
				}
				if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
					return "", command.NewExitStatusError("npm", &exec.ExitError{}, []string{"network error"})
				}
				return "", nil
			},
			expectError:   true,
			errorContains: "install Corepack",
		},
		{
			name: "corepack enable fails",
			runFunc: func() func(cmdStr string) (string, error) {
				callCount := 0
				return func(cmdStr string) (string, error) {
					if strings.Contains(cmdStr, "corepack") && strings.Contains(cmdStr, "--version") {
						callCount++
						if callCount == 1 {
							return "", fmt.Errorf("command not found")
						}
						return "0.17.0", nil
					}
					if strings.Contains(cmdStr, "npm") && strings.Contains(cmdStr, "install") {
						return "", nil
					}
					if strings.Contains(cmdStr, "enable") {
						return "", command.NewExitStatusError("corepack", &exec.ExitError{}, []string{"permission denied"})
					}
					return "", nil
				}
			}(),
			expectError:   true,
			errorContains: "enable Corepack",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.NewLogger(log.WithOutput(&logBuf))

			step := &Step{
				logger:     logger,
				cmdFactory: &mockFactory{runFunc: tt.runFunc},
			}

			err := step.ensureCorepack("test-dir")

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

			if tt.expectLogMsg != "" {
				logOutput := logBuf.String()
				if !strings.Contains(logOutput, tt.expectLogMsg) {
					t.Errorf("Expected log to contain %q, got %q", tt.expectLogMsg, logOutput)
				}
			}
		})
	}
}
