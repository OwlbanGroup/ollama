package llm

import (




)

// Request represents the structure of a request to the multimodal model.
type Request struct {
	ImageBase64 string `json:"image_base64"`
}

// CreateRequest constructs a request object from the base64 image string.
func CreateRequest(imageBase64 string) *Request {
	return &Request{
		ImageBase64: imageBase64,
	}
}

type Response struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// SendRequest sends the request to the multimodal model and returns the response.
func SendRequest(request *Request) (*Response, error) {
	// Mock response for demonstration purposes
	response := &Response{
		Content: "expected substring",
		Done:    true,
	}
	return response, nil
}
