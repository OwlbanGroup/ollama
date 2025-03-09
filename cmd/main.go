package cmd

import (
	"log"
	"time" // Added import for time package
	"github.com/ollama/ollama/common" // Added import for common package
	"github.com/spf13/cobra"
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
