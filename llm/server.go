package llm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/envconfig"
	"github.com/ollama/ollama/format"
	"github.com/ollama/ollama/gpu"
	"golang.org/x/sync/semaphore"
)

// CompletionRequest represents a request for model completion
type CompletionRequest struct {
	Prompt  string
	Options *api.Options
}

// CompletionResponse represents the response from a model completion request
type CompletionResponse struct {
	Content string
	Done    bool
}

func CreateRequest(imageBase64 string) CompletionRequest {
	return CompletionRequest{
		Prompt:  imageBase64,
		Options: nil, // Set options as needed
	}
}

func SendRequest(s *llmServer, req CompletionRequest) (CompletionResponse, error) {
	var response CompletionResponse
	err := s.Completion(context.Background(), req, func(resp CompletionResponse) {
		response = resp
	})
	return response, err
}

type LlamaServer interface {
	Ping(ctx context.Context) error
	WaitUntilRunning(ctx context.Context) error
	Completion(ctx context.Context, req CompletionRequest, fn func(CompletionResponse)) error
	Embed(ctx context.Context, input []string) ([][]float32, error)
	Detokenize(ctx context.Context, tokens []int) (string, error)

	Tokenize(ctx context.Context, content string) ([]int, error)
	Close() error
	EstimatedVRAM() uint64
	EstimatedTotal() uint64
	EstimatedVRAMByGPU(gpuID string) uint64
}

// llmServer is an instance of the llama.cpp server
type llmServer struct {
	port    int
	cmd     *exec.Cmd
	done    chan error
	status  *StatusWriter
	options api.Options

	estimate     MemoryEstimate
	totalLayers  uint64
	gpus         gpu.GpuInfoList
	loadDuration time.Duration
	loadProgress float32

	sem *semaphore.Weighted
}

// LoadModel will load a model from disk. The model must be in the GGML format.
func LoadModel(model string, maxArraySize int) (*GGML, error) {
	if _, err := os.Stat(model); err != nil {
		return nil, err
	}

	f, err := os.Open(model)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ggml, _, err := DecodeGGML(f, maxArraySize)
	return ggml, err
}

func Quantize(model string, fileType string) error {
	fmt.Printf("Quantizing model: %s of type: %s\n", model, fileType)
	return nil
}

// NewLlamaServer will run a server for the given GPUs
func NewLlamaServer(gpus gpu.GpuInfoList, model string, ggml *GGML, adapters, projectors []string, opts api.Options, numParallel int) (LlamaServer, error) {
	var err error
	var cpuRunner string
	var estimate MemoryEstimate
	var systemTotalMemory uint64
	var systemFreeMemory uint64

	totalMemory, freeMemory := gpu.GetCPUInfo()

	log.Printf("CPU Info - Total Memory: %d, Free Memory: %d", totalMemory, freeMemory)

	systemMemInfo := struct {
		TotalMemory uint64
		FreeMemory  uint64
	}{
		TotalMemory: totalMemory,
		FreeMemory:  freeMemory,
	}

	systemTotalMemory = systemMemInfo.TotalMemory
	systemFreeMemory = systemMemInfo.FreeMemory
	log.Println("system memory", "total", format.HumanBytes2(systemTotalMemory), "free", format.HumanBytes2(systemFreeMemory))

	if opts.NumGPU == 0 {
		gpus = gpu.GetCPUInfo()
	}
	if len(gpus) == 1 && gpus[0].Library == "cpu" {
		cpuRunner = serverForCpu()
		estimate = EstimateGPULayers(gpus, ggml, projectors, opts)
	} else {
		estimate = EstimateGPULayers(gpus, ggml, projectors, opts)

		switch {
		case gpus[0].Library == "metal" && estimate.VRAMSize > systemTotalMemory:
			opts.NumGPU = 0
		case gpus[0].Library != "metal" && estimate.Layers == 0:
			cpuRunner = serverForCpu()
			gpus = gpu.GetCPUInfo()
		case opts.NumGPU < 0 && estimate.Layers > 0 && gpus[0].Library != "cpu":
			opts.NumGPU = estimate.Layers
		}
	}

	vram, physicalMemory, freeMemory := gpu.GetGPUInfo()

	log.Printf("GPU Info - VRAM: %d, Physical Memory: %d, Free Memory: %d", vram, physicalMemory, freeMemory)

	log.Println("Retrieved GPU Info", "vram", vram, "physicalMemory", physicalMemory, "freeMemory", freeMemory)

	if runtime.GOOS == "linux" {
		systemMemoryRequired := estimate.TotalSize - estimate.VRAMSize
		available := systemFreeMemory
		if systemMemoryRequired > available {
			log.Println("model request too large for system", "requested", format.HumanBytes2(systemMemoryRequired), "available", available, "total", format.HumanBytes2(systemTotalMemory), "free", format.HumanBytes2(systemFreeMemory))
			return nil, fmt.Errorf("model requires more system memory (%s) than is available (%s)", format.HumanBytes2(systemMemoryRequired), format.HumanBytes2(available))
		}
	}

	estimate.log()

	finalErr := errors.New("no suitable llama servers found")

	if len(adapters) > 1 {
		return nil, errors.New("ollama supports only one lora adapter, but multiple were provided")
	}

	availableServers := getAvailableServers()
	if len(availableServers) == 0 {
		if runtime.GOOS != "windows" {
			log.Println("llama server binary disappeared, reinitializing payloads")
			err = Init()
			if err != nil {
				log.Println("failed to reinitialize payloads", "error", err)
				return nil, err
			}
			availableServers = getAvailableServers()
		} else {
			return nil, finalErr
		}
	}
	var servers []string
	if cpuRunner != "" {
		servers = []string{cpuRunner}
	} else {
		servers = serversForGpu(gpus[0])
	}
	demandLib := envconfig.LLMLibrary
	if demandLib != "" {
		serverPath := availableServers[demandLib]
		if serverPath == "" {
			log.Println(fmt.Sprintf("Invalid OLLAMA_LLM_LIBRARY %s - not found", demandLib))
		} else {
			log.Println("user override", "OLLAMA_LLM_LIBRARY", demandLib, "path", serverPath)
			servers = []string{demandLib}
			if strings.HasPrefix(demandLib, "cpu") {
				opts.NumGPU = -1
			}
		}
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers found for %v", gpus)
	}

	params := []string{
		"--model", model,
		"--ctx-size", fmt.Sprintf("%d", opts.NumCtx),
		"--batch-size", fmt.Sprintf("%d", opts.NumBatch),
		"--embedding",
	}

	_ = params

	return &llmServer{
		port:    8080,
		options: opts,
	}, nil
}

func (s *llmServer) Completion(ctx context.Context, req CompletionRequest, fn func(CompletionResponse)) error {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	// Here you would implement the actual logic to process the request
	response := CompletionResponse{
		Content: fmt.Sprintf("Processed prompt: %s", req.Prompt),
		Done:    true,
	}
	fn(response)
	return nil
}


func (s *llmServer) Embed(ctx context.Context, input []string) ([][]float32, error) {
	return [][]float32{{0.0}}, nil
}

func (s *llmServer) EstimatedTotal() uint64 {
	return 0
}

func (s *llmServer) EstimatedVRAM() uint64 {
	return 0
}

func (s *llmServer) EstimatedVRAMByGPU(gpuID string) uint64 {
	return 0
}

func (s *llmServer) Ping(ctx context.Context) error {
	return nil
}

func (s *llmServer) Tokenize(ctx context.Context, content string) ([]int, error) {
	return []int{0}, nil
}

func (s *llmServer) WaitUntilRunning(ctx context.Context) error {
	return nil
}

func (s *llmServer) Detokenize(ctx context.Context, tokens []int) (string, error) {
	return "detokenized string", nil
}

func (s *llmServer) Close() error {
	if s.cmd != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}
