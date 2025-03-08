package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time" // Importing the time package

	"github.com/gin-gonic/gin" // Ensure gin is imported

	"github.com/ollama/ollama/api" // Ensure api is imported
	"github.com/ollama/ollama/llm"
)

var mode string = gin.DebugMode

type Server struct {
	addr  net.Addr
	sched *Scheduler
}

func init() {
	switch mode {
	case gin.DebugMode:
	case gin.ReleaseMode:
	case gin.TestMode:
	default:
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
}

var errRequired = errors.New("is required")
var errBadTemplate = errors.New("template error")

func modelOptions(model *Model, requestOpts map[string]interface{}) (api.Options, error) {
	opts := api.DefaultOptions()
	if err := opts.FromMap(model.Options); err != nil {
		return api.Options{}, err
	}

	if err := opts.FromMap(requestOpts); err != nil {
		return api.Options{}, err
	}

	return opts, nil
}

// scheduleRunner schedules a runner after validating inputs such as capabilities and model options.
// It returns the allocated runner, model instance, and consolidated options if successful and error otherwise.
func (s *Server) scheduleRunner(ctx context.Context, name string, caps []Capability, requestOpts map[string]any, keepAlive *api.Duration) (llm.LlamaServer, *Model, *api.Options, error) {
	if name == "" {
		return llm.LlamaServer{}, nil, nil, fmt.Errorf("model %w", errRequired)
	}

	model, err := GetModel(name)
	if err != nil {
		return llm.LlamaServer{}, nil, nil, err
	}

	if err := model.CheckCapabilities(caps...); err != nil {
		return llm.LlamaServer{}, nil, nil, fmt.Errorf("%s %w", name, err)
	}

	opts, err := modelOptions(model, requestOpts)
	if err != nil {
		return llm.LlamaServer{}, nil, nil, err
	}

	runnerCh, errCh := s.sched.GetRunner(ctx, model, opts, keepAlive)
	var runner *runnerRef
	select {
	case runner = <-runnerCh:
	case err = <-errCh:
		return llm.LlamaServer{}, nil, nil, err
	}

	return runner.llama, model, &opts, nil
}

func (s *Server) GenerateHandler(c *gin.Context) {
	var req api.GenerateRequest
	if err := c.ShouldBindJSON(&req); errors.Is(err, io.EOF) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format != "" && req.Format != "json" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "format must be empty or \"json\""})
		return
	} else if req.Raw && (req.Template != "" || req.System != "" || len(req.Context) > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "raw mode does not support template, system, or context"})
		return
	}

	caps := []Capability{CapabilityCompletion}
	if req.Suffix != "" {
		caps = append(caps, CapabilityInsert)
	}

	r, opts, err := s.scheduleRunner(c.Request.Context(), req.Model, caps, req.Options, req.KeepAlive) // Fixing assignment mismatch



	if errors.Is(err, errCapabilityCompletion) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%q does not support generate", req.Model)})
		return
	} else if err != nil {
		handleScheduleError(c, req.Model, err)
		return
	}

	if req.Prompt == "" {
		c.JSON(http.StatusOK, api.GenerateResponse{
			Model:     req.Model,
			CreatedAt: time.Now().UTC(),
			Done:      true,
		})
		return
	}

	prompt := req.Prompt
	response, err := r.Completion(c.Request.Context(), llm.CompletionRequest{
		Prompt:     prompt,
		Temperature: opts.Temperature,
		MaxTokens:  opts.MaxTokens,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func Serve(ln net.Listener) error { // Ensure s is passed as a parameter

	router := gin.Default()

	// Register routes
	router.POST("/generate", s.GenerateHandler) // Ensure the correct reference to the method

	// Start the server
	return s.Serve(ln) // Ensure s is defined and initialized

}

// Define the missing functions here
func handleScheduleError(c *gin.Context, model string, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error scheduling model %s: %v", model, err)})
}

func streamResponse(c *gin.Context, ch chan any) {
	for rr := range ch {
		switch t := rr.(type) {
		case api.GenerateResponse:
			c.JSON(http.StatusOK, t)
		case gin.H:
			c.JSON(http.StatusInternalServerError, t)
		}
	}
}
