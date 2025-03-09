package server

import (
	"log"
	"time" // Added import for time package
	"github.com/ollama/ollama/gpu"
	"github.com/ollama/ollama/llm"
	"github.com/ollama/ollama/common" // Added import for common package
)





// Scheduler manages the scheduling of tasks
type Scheduler struct {
	llama *llm.LlamaServer
	pendingReqCh  chan *llm.LlmRequest // Updated to use common.runnerRef
	finishedReqCh chan *llm.LlmRequest
	expiredCh     chan *common.runnerRef // Updated to use common.runnerRef
	unloadedCh    chan any
	loaded        map[string]*common.runnerRef // Updated to use common.runnerRef
	newServerFn   func(gpu.GpuInfoList, string, *llm.GGML, []string, []string, api.Options, int) (llm.LlamaServer, error)
	getGpuFn      func() (gpu.GpuInfoList, error)
	getCpuFn      func() (interface{}, error) // Updated to use interface{} for cpuInfo
	reschedDelay  time.Duration
	loadFn        func(req *llm.LlmRequest, ggml *llm.GGML, gpus gpu.GpuInfoList, numParallel int)
}



// NewScheduler creates a new Scheduler instance
func NewScheduler() *Scheduler {
	llamaServer := llm.NewLlamaServer("localhost", 8080)
	cpuTotal, cpuFree := gpu.GetCPUInfo()
	vram, gpuTotal, gpuFree := gpu.GetGPUInfo()

	log.Printf("CPU Info - Total: %d, Free: %d", cpuTotal, cpuFree)
	log.Printf("GPU Info - VRAM: %d, Total: %d, Free: %d", vram, gpuTotal, gpuFree)

	return &Scheduler{
		llama: llamaServer,
	}
}

// Start begins the scheduling process
func (s *Scheduler) Start() {
	if s.llama != nil {
		log.Println("Llama server is running.")
	} else {
		log.Println("Llama server is not initialized.")
	}
}

// Other methods for Scheduler can be added here
