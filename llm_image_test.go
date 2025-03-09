package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ollama/ollama/llm" // Import the llm package using the full module path
)

func TestIntegrationMultimodal(t *testing.T) {
	// Decode the base64 image and prepare the request
	imageBase64 := "..." // Replace with actual base64 image string
	request := llm.CreateRequest(imageBase64)

	// Send the request to the multimodal model
	response, err := llm.SendRequest(request)
	require.NoError(t, err)

	// Check the response
	assert.NotNil(t, response)
	assert.Contains(t, response.Content, "expected substring") // Adjust based on expected response
	assert.Equal(t, true, response.Done) // Ensure the response indicates completion
}
