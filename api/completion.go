package api

// Removed import of common to eliminate potential circular dependency

type CompletionRequest struct {
	Prompt     string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens  int     `json:"max_tokens"`
}

type CompletionResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
