//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ollama/ollama/api"
)

func TestContextExhaustion(t *testing.T) {
	startTime := time.Now() // Start time for logging duration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		duration := time.Since(startTime)
		t.Logf("TestContextExhaustion completed in %v", duration) // Log the duration of the test
	}()

	defer cancel()
	// Set up the test data
	req := api.GenerateRequest{
		Model:  "llama2",
		Prompt: "Write me a story with a ton of emojis?",
		Stream: &stream,
		Options: map[string]interface{}{
			"temperature": 0,
			"seed":        123,
			"num_ctx":     128,
		},
	}
	client, _, cleanup := InitServerConnection(ctx, t)
	defer cleanup()
	if err := PullIfMissing(ctx, client, req.Model); err != nil {
		t.Fatalf("PullIfMissing failed: %v", err)
	}
	DoGenerate(ctx, t, client, req, []string{"once", "upon", "lived"}, 120*time.Second, 10*time.Second)
}
