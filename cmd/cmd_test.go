package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestCommandExecution(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
	}{
		{
			name:    "Test help command",
			command: "help",
			args:    []string{},
			want:    "Usage:",
		},
		{
			name:    "Test version command",
			command: "version",
			args:    []string{},
			want:    "version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(tt.command, tt.args...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			if got := out.String(); got != tt.want {
				t.Errorf("Command output = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Setup code if needed
	code := m.Run()
	// Teardown code if needed
	os.Exit(code)
}
