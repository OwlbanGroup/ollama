package llm

// #cgo CFLAGS: -Illama.cpp -Illama.cpp/include -Illama.cpp/ggml/include
// #cgo LDFLAGS: -lllama -lggml -lstdc++ -lpthread
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/build/darwin/arm64_static -L${SRCDIR}/build/darwin/arm64_static/src -L${SRCDIR}/build/darwin/arm64_static/ggml/src -framework Accelerate -framework Metal
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/build/darwin/x86_64_static -L${SRCDIR}/build/darwin/x86_64_static/src -L${SRCDIR}/build/darwin/x86_64_static/ggml/src
// #cgo windows,amd64 LDFLAGS: -static-libstdc++ -static-libgcc -static -L${SRCDIR}/build/windows/amd64_static -L${SRCDIR}/build/windows/amd64_static/src -L${SRCDIR}/build/windows/amd64_static/ggml/src
// #cgo windows,arm64 LDFLAGS: -static-libstdc++ -static-libgcc -static -L${SRCDIR}/build/windows/arm64_static -L${SRCDIR}/build/windows/arm64_static/src -L${SRCDIR}/build/windows/arm64_static/ggml/src
// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/build/linux/x86_64_static -L${SRCDIR}/build/linux/x86_64_static/src -L${SRCDIR}/build/linux/x86_64_static/ggml/src
// #cgo linux,arm64 LDFLAGS: -L${SRCDIR}/build/linux/arm64_static -L${SRCDIR}/build/linux/arm64_static/src -L${SRCDIR}/build/linux/arm64_static/ggml/src
// #include <stdlib.h>
// #include "llama.h"
import "C"
import (
	"fmt"
	"unsafe"
	"github.com/ollama/ollama/common" // Added import for common package
)

func LoadModel(modelPath string) error {
	model := C.llama_load_model(C.CString(modelPath))
	if model == nil {
		return fmt.Errorf("failed to load model from path: %s", modelPath)
	}

	fmt.Printf("Loading model from: %s\n", modelPath)
	return nil
}

// CompletionResponse represents the response structure for completions.
type CompletionResponse struct {
	DoneReason         string  `json:"done_reason"`
	PromptEvalCount    int     `json:"prompt_eval_count"`
	PromptEvalDuration float64 `json:"prompt_eval_duration"`
	EvalCount          int     `json:"eval_count"`
	EvalDuration       float64 `json:"eval_duration"`
}

// ImageData represents image-related data.
type ImageData struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// CompletionRequest represents the request structure for completions.
type CompletionRequest struct {
	Prompt     string     `json:"prompt"`
	MaxTokens  int        `json:"max_tokens"`
	Temperature float64    `json:"temperature"`
	TopP       float64    `json:"top_p"`
	Images     []ImageData `json:"images,omitempty"`
	Format     string     `json:"format,omitempty"`
}

func SystemInfo() string {
	return C.GoString(C.llama_print_system_info())
}

func Quantize(infile, outfile string, ftype fileType) error {
	cinfile := C.CString(infile)
	defer C.free(unsafe.Pointer(cinfile))

	coutfile := C.CString(outfile)
	defer C.free(unsafe.Pointer(coutfile))

	params := C.llama_model_quantize_default_params()
	params.nthread = -1
	params.ftype = ftype.Value()

	if rc := C.llama_model_quantize(cinfile, coutfile, &params); rc != 0 {
		return fmt.Errorf("failed to quantize model. This model architecture may not be supported, or you may need to upgrade Ollama to the latest version")
	}

	return nil
}
