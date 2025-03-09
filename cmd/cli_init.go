package cmd

import (
	"github.com/spf13/cobra"
)

// NewCLI initializes the command-line interface for the application.
func NewCLI() *cobra.Command {
	// Create a new command
	cmd := &cobra.Command{
		Use:   "ollama",
		Short: "Ollama CLI",
		Long:  "A command-line interface for the Ollama application.",
	}

	// Define subcommands here
	// cmd.AddCommand(...)

	return cmd
}
