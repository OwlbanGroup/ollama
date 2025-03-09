package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// createRequest is a utility function to create an HTTP request for testing
func createRequest(t *testing.T, handler http.HandlerFunc, reqBody interface{}) *httptest.ResponseRecorder {
	t.Helper() // Added helper to mark this function as a test helper

	var body *bytes.Buffer
	if reqBody != nil {
		body = new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(reqBody)
		if err != nil {
			t.Fatalf("failed to encode request body: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// createBinFile is a utility function to create a binary file for testing
func createBinFile(t *testing.T, data interface{}) string {
	t.Helper()
	// Logic to create a binary file and return its name
	return "dummy.bin" // Placeholder implementation
}
