package pkgmanager

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bitrise-io/go-utils/v2/log"
)

var pkgFieldPresent = `{
	"name": "test-project",
	"version": "1.0.0",
	"packageManager": "yarn@4.1.0"
}`

var pkgFieldMissing = `{
	"name": "test-project",
	"version": "1.0.0"
}`

func TestPrintPackageManagerInfo(t *testing.T) {
	tests := []struct {
		name           string
		packageJSON    string
		expectedOutput string
	}{
		{
			name:           "packageManager field present",
			packageJSON:    pkgFieldPresent,
			expectedOutput: "package.json specifies packageManager: yarn@4.1.0",
		},
		{
			name:           "packageManager field missing",
			packageJSON:    pkgFieldMissing,
			expectedOutput: "No packageManager field found in package.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.NewLogger(log.WithOutput(&buf))
			mockFS := fstest.MapFS{
				"package.json": {
					Data: []byte(tt.packageJSON),
				},
			}
			step := &Step{
				logger: logger,
				fs:     mockFS,
			}

			step.printPackageManagerInfo()

			output := buf.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestPrintYarnrcInfo(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		expectedOutput string
	}{
		{
			name: "yarnrc file found",
			files: map[string]string{
				".yarnrc": "registry \"https://registry.npmjs.org/\"",
			},
			expectedOutput: ".yarnrc file found (Yarn 1.x config)",
		},
		{
			name: "yarnrc.yml file found",
			files: map[string]string{
				".yarnrc.yml": "nodeLinker: node-modules",
			},
			expectedOutput: ".yarnrc.yml file found (Yarn 2+ config)",
		},
		{
			name:           "no yarnrc files",
			files:          map[string]string{},
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.NewLogger(log.WithOutput(&buf))
			mockFS := fstest.MapFS{}
			for name, content := range tt.files {
				mockFS[name] = &fstest.MapFile{
					Data: []byte(content),
				}
			}
			step := &Step{
				logger: logger,
				fs:     mockFS,
			}

			step.printYarnrcInfo()

			output := buf.String()
			if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
