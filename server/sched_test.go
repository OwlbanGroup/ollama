package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/app/lifecycle"
	"github.com/ollama/ollama/gpu"
	"github.com/ollama/ollama/llm"
	"github.com/stretchr/testify/require"
)



func init() {
	os.Setenv("OLLAMA_DEBUG", "1")
	lifecycle.InitLogging()
}

type mockLlm struct {
	estimatedVRAM uint64
}

func (m *mockLlm) EstimatedVRAM() uint64 {
	return m.estimatedVRAM
}

func TestInitScheduler(t *testing.T) {
	s := InitScheduler()
	s.loadedMu.Lock()
	require.NotNil(t, s.loaded)
	s.loadedMu.Unlock()
}

func TestLoad(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer done()
	s := InitScheduler()
	req := &llm.LlmRequest{
		Model: "foo",
	}
	s.newServerFn = func(gpus gpu.GpuInfoList, model string, ggml *llm.GGML, adapters []string, projectors []string, opts api.Options, numParallel int) (llm.LlamaServer, error) {
		return &mockLlm{estimatedVRAM: 10}, nil
	}
	s.load(req, nil, nil, 0)
	require.Empty(t, req.successCh)
	require.Len(t, req.errCh, 1)
}

func TestRequestsSameModelSameRequest(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer done()
	s := InitScheduler()
	a := &llm.LlmRequest{
		Model: "ollama-model-1",
	}
	b := &llm.LlmRequest{
		Model: "ollama-model-1",
	}

	s.newServerFn = func(gpus gpu.GpuInfoList, model string, ggml *llm.GGML, adapters []string, projectors []string, opts api.Options, numParallel int) (llm.LlamaServer, error) {
		return &mockLlm{estimatedVRAM: 10}, nil
	}
	s.load(a, nil, nil, 0)
	s.load(b, nil, nil, 0)

	require.Equal(t, a.Model, b.Model)
}

func TestUnloadAllRunners(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer done()

	s := InitScheduler()
	r1 := &runnerRef{llama: &mockLlm{estimatedVRAM: 10}, numParallel: 1}
	r2 := &runnerRef{llama: &mockLlm{estimatedVRAM: 20}, numParallel: 1}

	s.loadedMu.Lock()
	s.loaded["a"] = r1
	s.loaded["b"] = r2
	s.loadedMu.Unlock()
	s.unloadAllRunners()

	require.True(t, r1.llama.(*mockLlm).closeCalled)
	require.True(t, r2.llama.(*mockLlm).closeCalled)
}
