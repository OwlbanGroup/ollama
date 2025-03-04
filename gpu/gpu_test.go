package gpu

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicGetGPUInfo(t *testing.T) {
	vram, physicalMemory, freeMemory := GetGPUInfo()
	assert.Greater(t, vram, uint64(0))
	assert.Greater(t, physicalMemory, uint64(0))
	assert.Greater(t, freeMemory, uint64(0))

}

func TestCPUMemInfo(t *testing.T) {
	info, err := GetCPUMem()
	require.NoError(t, err)
	switch runtime.GOOS {
	case "darwin":
		t.Skip("CPU memory not populated on darwin")
	case "linux", "windows":
		assert.Greater(t, info.TotalMemory, uint64(0))
		assert.Greater(t, info.FreeMemory, uint64(0))
	default:
		return
	}
}

// TODO - add some logic to figure out card type through other means and actually verify we got back what we expected
