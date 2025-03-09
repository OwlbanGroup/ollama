package main

import (
	"context"
	"github.com/ollama/ollama/cmd" // Ensure the correct import path
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewCLI().ExecuteContext(context.Background()))
}
