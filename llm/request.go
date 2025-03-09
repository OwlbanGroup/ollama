package llm

type LlmRequest struct {
	Model string
	// Add other necessary fields
}

type CompletionRequest struct {
	Prompt string
	// Add other necessary fields
	// Removed unused import
}
