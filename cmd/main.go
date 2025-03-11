package cmd

import (
	"log"

	"github.com/ollama/ollama/common" // Added import for common package
	"context" // Importing context package

	"github.com/spf13/cobra"
	// Removed import of cmd package


)



func main() {

	// Initialize the Llama server
	llamaServer := common.LlamaServer{
		Address: "localhost",
		Port:    8080,
	}

	log.Printf("Llama server initialized: %+v", llamaServer)
	cobra.CheckErr(cmd.NewCLI().ExecuteContext(context.Background()))
}
