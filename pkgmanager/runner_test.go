package pkgmanager

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

func TestExecuteYarnCommand(t *testing.T) {
	tests := []struct {
		name          string
		runFunc       func(cmdStr string) (string, error)
		expectError   bool
		errorContains string
		expectLogMsg  string
	}{
		{
			name: "yarn command succeeds",
			runFunc: func(cmdStr string) (string, error) {
				return "", nil
			},
			expectError: false,
		},
		{
			name: "yarn command fails",
			runFunc: func(cmdStr string) (string, error) {
				return "Build failed", command.NewExitStatusError("yarn", &exec.ExitError{}, []string{"Build failed"})
			},
			expectError:   true,
			errorContains: "yarn command failed",
		},
		{
			name: "network error triggers warning",
			runFunc: func(cmdStr string) (string, error) {
				return "There appears to be trouble with your network connection. Retrying...", command.NewExitStatusError("yarn", &exec.ExitError{}, []string{"Network error"})
			},
			expectError:   true,
			errorContains: "yarn command failed",
			expectLogMsg:  "network issues",
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

			err := step.ExecuteYarnCommand("test-dir", []string{}, []string{"test-arg"})

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
