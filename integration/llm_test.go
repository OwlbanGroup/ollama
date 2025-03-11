//go:build integration

package integration_tests

import (
	"context"
	"testing"
	"time"
)

type GenerateRequest struct {
	Model  string
	Prompt string
	Stream *bool
	Options map[string]interface{}
}

// Additional test functions can be added here


var req = [2]GenerateRequest{
	{
		Model:  "orca-mini",
		Prompt: "why is the ocean blue?",
		Stream: new(bool),
		Options: map[string]interface{}{
			"seed":        42,
			"temperature": 0.0,
		},
	}, {
		Model:  "orca-mini",
		Prompt: "what is the origin of the us thanksgiving holiday?",
		Stream: new(bool),
		Options: map[string]interface{}{
			"seed":        42,
			"temperature": 0.0,
		},
	},
}


// TODO - this would ideally be in the llm package, but that would require some refactoring of interfaces in the server
//        package to avoid circular dependencies

var resp = [2][]string{
	{"sunlight"},
	{"england", "english", "massachusetts", "pilgrims"},
}


func TestIntegrationSimpleOrcaMini(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	GenerateTestHelper(ctx, t, req[0], resp[0])
}
